package handler

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kotche/url-shortening-service/internal/app/service"
	"github.com/kotche/url-shortening-service/internal/config"
)

type Storage interface {
	Add(url *service.URL)
	GetByID(id string) (*service.URL, error)
}

type Handler struct {
	st     Storage
	router *chi.Mux
}

func NewHandler(st Storage) *Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	handler := &Handler{
		st:     st,
		router: router,
	}

	handler.setRouting()

	return handler
}

func (h *Handler) setRouting() {
	h.router.Get("/{id}", h.handleGet)
	h.router.Post("/", h.handlePost)
}

func (h *Handler) GetRouter() *chi.Mux {
	return h.router
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
		urlModel, _ := h.st.GetByID(shortURL)

		if urlModel == nil {
			urlModel = service.NewURL(originURL, shortURL)
			h.st.Add(urlModel)
		} else if urlModel.GetOriginal() != originURL {
			continue
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(config.ServerAddrForURL + urlModel.GetShort()))
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
	w.Header().Set("Location", url.GetOriginal())
	w.WriteHeader(http.StatusTemporaryRedirect)
}
