package handler

import (
	"github.com/kotche/url-shortening-service/internal/app/service"
	"github.com/kotche/url-shortening-service/internal/app/storage"
	"github.com/kotche/url-shortening-service/internal/config"
	"io"
	"net/http"
	"strings"
)

type Handler struct {
	st storage.Storage
}

func NewHandler(st storage.Storage) *Handler {
	return &Handler{st: st}
}

func (h *Handler) Handlers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.handlePost(w, r)
		case http.MethodGet:
			h.handleGet(w, r)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path

	if len(path) > 0 && path != "/" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
			w.Write([]byte(err.Error()))
			return
		}

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

	urlParts := strings.Split(r.URL.Path, "/")

	if len(urlParts) == 2 && urlParts[1] != "" {

		url, err := h.st.GetByID(urlParts[1])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Location", url.GetOriginal())
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
