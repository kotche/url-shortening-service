package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kotche/url-shortening-service/internal/app/middlewares"
	"github.com/kotche/url-shortening-service/internal/app/service"
	"github.com/kotche/url-shortening-service/internal/app/usecase"
	"github.com/kotche/url-shortening-service/internal/config"
)

type ICookieManager interface {
	GetUserID(r *http.Request) string
}

type Handler struct {
	Service *service.Service
	Router  *chi.Mux
	Conf    *config.Config
	Cm      ICookieManager
}

func NewHandler(service *service.Service, conf *config.Config) *Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middlewares.GzipHandle)
	router.Use(middlewares.UserCookieHandle)

	handler := &Handler{
		Service: service,
		Router:  router,
		Conf:    conf,
		Cm:      CookieManager{},
	}

	handler.setRouting()

	return handler
}

func (h *Handler) setRouting() {
	h.Router.Get("/{id}", h.handleGet)
	h.Router.Post("/", h.handlePost)
	h.Router.Post("/api/shorten", h.handlePostJSON)
	h.Router.Get("/api/user/urls", h.handleGetUserURLs)
	h.Router.Get("/ping", h.handlePing)
	h.Router.Post("/api/shorten/batch", h.handlePostShortenBatch)
	h.Router.Delete("/api/user/urls", h.handleDeleteURLs)
}

func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	urlRead, err := io.ReadAll(r.Body)
	if err != nil || len(urlRead) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originURL := string(urlRead)
	userID := h.Cm.GetUserID(r)

	status := http.StatusCreated
	var shortenURL string

	urlModel, err := h.Service.GetURLModel(userID, originURL)

	if errors.As(err, &usecase.ConflictURLError{}) {
		e := err.(usecase.ConflictURLError)
		status = http.StatusConflict
		shortenURL = e.ShortenURL
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		shortenURL = urlModel.Short
	}

	w.WriteHeader(status)
	_, _ = w.Write([]byte(h.Conf.BaseURL + "/" + shortenURL))
}

func (h *Handler) handlePostJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid content type", http.StatusBadRequest)
		return
	}

	originURLReceiver := &struct {
		OriginURL string `json:"url"`
	}{}

	err := json.NewDecoder(r.Body).Decode(originURLReceiver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	originURL := originURLReceiver.OriginURL

	if originURL == "" {
		http.Error(w, "URL is empty", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	userID := h.Cm.GetUserID(r)

	status := http.StatusCreated
	var shortenURL string

	urlModel, err := h.Service.GetURLModel(userID, originURL)

	if errors.As(err, &usecase.ConflictURLError{}) {
		e := err.(usecase.ConflictURLError)
		status = http.StatusConflict
		shortenURL = e.ShortenURL
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		shortenURL = urlModel.Short
	}

	shortURLSender := &struct {
		ShortURL string `json:"result"`
	}{
		ShortURL: h.Conf.BaseURL + "/" + shortenURL,
	}

	response, err := json.Marshal(shortURLSender)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	_, _ = w.Write(response)
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "id")

	url, err := h.Service.GetURLModelByID(shortURL)

	if errors.As(err, &usecase.GoneError{}) {
		w.WriteHeader(http.StatusGone)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url.Origin)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) handleGetUserURLs(w http.ResponseWriter, r *http.Request) {

	userID := h.Cm.GetUserID(r)

	if userID == "" {
		http.Error(w, "user ID is empty", http.StatusInternalServerError)
		return
	}

	userUrls, err := h.Service.GetUserURLs(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(userUrls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	type Output struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	outputList := make([]Output, 0, len(userUrls))

	for _, v := range userUrls {
		p := Output{
			ShortURL:    h.Conf.BaseURL + "/" + v.Short,
			OriginalURL: v.Origin,
		}
		outputList = append(outputList, p)
	}

	userUrlsJSON, err := json.Marshal(outputList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(userUrlsJSON)
}

func (h *Handler) handlePing(w http.ResponseWriter, _ *http.Request) {
	err := h.Service.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handlePostShortenBatch(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputDataList := make([]usecase.InputCorrelationURL, 0)
	err = json.Unmarshal(body, &inputDataList)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID := h.Cm.GetUserID(r)

	ctx := context.Background()

	outputDataList, err := h.Service.ShortenBatch(ctx, userID, inputDataList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for ind := range outputDataList {
		outputDataList[ind].Short = h.Conf.BaseURL + "/" + outputDataList[ind].Short
	}

	correlationURLs, _ := json.Marshal(outputDataList)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(correlationURLs)
}

func (h *Handler) handleDeleteURLs(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	toDelete := make([]string, 0)
	err = json.Unmarshal(body, &toDelete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID := h.Cm.GetUserID(r)

	go func() {
		h.Service.DeleteURLs(userID, toDelete)
	}()

	w.WriteHeader(http.StatusAccepted)
}
