package handler

import (
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

type Handler struct {
	service *service.Service
	router  *chi.Mux
	conf    *config.Config
}

func (h *Handler) GetRouter() *chi.Mux {
	return h.router
}

func NewHandler(service *service.Service, conf *config.Config) *Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middlewares.GzipHandle)
	router.Use(middlewares.UserCookieHandle)

	handler := &Handler{
		service: service,
		router:  router,
		conf:    conf,
	}

	handler.setRouting()

	return handler
}

func (h *Handler) setRouting() {
	h.router.Get("/{id}", h.handleGet)
	h.router.Post("/", h.handlePost)
	h.router.Post("/api/shorten", h.handlePostJSON)
	h.router.Get("/api/user/urls", h.handleGetUserURLs)
	h.router.Get("/ping", h.handlePing)
	h.router.Post("/api/shorten/batch", h.handlePostShortenBatch)
}

func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	urlRead, err := io.ReadAll(r.Body)
	if err != nil || len(urlRead) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originURL := string(urlRead)
	userID := r.Context().Value(config.UserIDCookieName).(string)

	urlModel, err := h.service.GetURLModel(userID, originURL)
	if err != nil {
		if errors.As(err, &usecase.ErrConflictURL{}) {
			e := err.(usecase.ErrConflictURL)
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(h.conf.BaseURL + "/" + e.ShortenURL))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.conf.BaseURL + "/" + urlModel.Short))
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

	userID := r.Context().Value(config.UserIDCookieName).(string)

	var shortenURL string

	urlModel, err := h.service.GetURLModel(userID, originURL)
	if err != nil {
		if errors.As(err, &usecase.ErrConflictURL{}) {
			e := err.(usecase.ErrConflictURL)
			w.WriteHeader(http.StatusConflict)
			shortenURL = e.ShortenURL
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		w.WriteHeader(http.StatusCreated)
		shortenURL = urlModel.Short
	}

	shortURLSender := &struct {
		ShortURL string `json:"result"`
	}{
		ShortURL: h.conf.BaseURL + "/" + shortenURL,
	}

	response, err := json.Marshal(shortURLSender)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "id")

	url, err := h.service.GetURLModelByID(shortURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Location", url.Origin)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) handleGetUserURLs(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(config.UserIDCookieName).(string)

	if userID == "" {
		http.Error(w, "user ID is empty", http.StatusInternalServerError)
		return
	}

	userUrls, err := h.service.GetUserURLs(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
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
			ShortURL:    h.conf.BaseURL + "/" + v.Short,
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

func (h *Handler) handlePing(w http.ResponseWriter, r *http.Request) {
	err := h.service.Ping()
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

	inputDataList := make([]service.InputCorrelationURL, 0)
	err = json.Unmarshal(body, &inputDataList)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value(config.UserIDCookieName).(string)
	outputDataList, err := h.service.ShortenBatch(userID, inputDataList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for ind := range outputDataList {
		outputDataList[ind].Short = h.conf.BaseURL + "/" + outputDataList[ind].Short
	}

	correlationURLs, _ := json.Marshal(outputDataList)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(correlationURLs)
}
