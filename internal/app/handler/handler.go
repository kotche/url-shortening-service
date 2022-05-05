package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kotche/url-shortening-service/internal/app/service"
	"github.com/kotche/url-shortening-service/internal/config"
)

type Storage interface {
	Add(url *service.URL) error
	GetByID(id string) (*service.URL, error)
	Close() error
}

type Handler struct {
	st     Storage
	router *chi.Mux
	conf   *config.Config
}

func (h *Handler) GetRouter() *chi.Mux {
	return h.router
}

func NewHandler(st Storage, conf *config.Config) *Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	handler := &Handler{
		st:     st,
		router: router,
		conf:   conf,
	}

	handler.setRouting()

	return handler
}

func (h *Handler) setRouting() {
	h.router.Get("/{id}", h.handleGet)
	h.router.Post("/", h.handlePost)
	h.router.Post("/api/shorten", h.handlePostJSON)
}

func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	urlRead, err := io.ReadAll(r.Body)
	if err != nil || len(urlRead) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originURL := string(urlRead)

	for {
		shortURL := service.MakeShortURL()
		urlModel, err := h.st.GetByID(shortURL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if urlModel == nil {
			urlModel = service.NewURL(originURL, shortURL)
			err := h.st.Add(urlModel)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else if urlModel.Origin != originURL {
			continue
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(h.conf.BaseURL + "/" + urlModel.Short))
		break
	}
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

	for {
		shortURL := service.MakeShortURL()
		urlModel, _ := h.st.GetByID(shortURL)

		if urlModel == nil {
			urlModel = service.NewURL(originURL, shortURL)
			err := h.st.Add(urlModel)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else if urlModel.Origin != originURL {
			continue
		}

		shortURLSender := &struct {
			ShortURL string `json:"result"`
		}{
			ShortURL: h.conf.BaseURL + "/" + urlModel.Short,
		}

		response, err := json.Marshal(shortURLSender)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
		break
	}
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "id")

	url, err := h.st.GetByID(shortURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Location", url.Origin)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
