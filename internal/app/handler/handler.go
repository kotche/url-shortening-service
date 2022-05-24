package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kotche/url-shortening-service/internal/app/middlewares"
	"github.com/kotche/url-shortening-service/internal/app/service"
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
}

func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	urlRead, err := io.ReadAll(r.Body)
	if err != nil || len(urlRead) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originURL := string(urlRead)
	userID := r.Context().Value(config.UserIDCookie).(string)

	urlModel, err := h.service.GetURLModel(userID, originURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

	userID := r.Context().Value(config.UserIDCookie).(string)

	urlModel, err := h.service.GetURLModel(userID, originURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURLSender := &struct {
		ShortURL string `json:"result"`
	}{
		ShortURL: h.conf.BaseURL + "/" + urlModel.Short,
	}

	response, err := json.Marshal(shortURLSender)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
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

	userID := r.Context().Value(config.UserIDCookie).(string)

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

	userUrlsJSON, err := json.Marshal(userUrls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(userUrlsJSON)
	w.Header().Add("Content-Type", "application/json")
}
