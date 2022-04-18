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
			w.WriteHeader(400)
		}
	}
}

func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path

	if len(path) > 0 && path != "/" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	urlRead, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originUrl := string(urlRead)

	for {
		shortUrl := service.MakeShortURl()
		urlModel, _ := h.st.GetById(shortUrl)

		if urlModel == nil {
			urlModel = service.NewUrl(originUrl, shortUrl)
			h.st.Add(urlModel)
		} else if urlModel != nil && urlModel.GetOriginal() != originUrl {
			continue
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(config.ServerAddrForUrl + urlModel.GetShort()))
		break
	}
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {

	urlParts := strings.Split(r.URL.Path, "/")

	if len(urlParts) == 2 && urlParts[1] != "" {

		url, err := h.st.GetById(urlParts[1])
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", url.GetOriginal())
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
