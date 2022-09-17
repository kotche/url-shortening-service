package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kotche/url-shortening-service/internal/app/config"
	"github.com/kotche/url-shortening-service/internal/app/model"
	"github.com/kotche/url-shortening-service/internal/app/service"
	middlewares2 "github.com/kotche/url-shortening-service/internal/app/transport/rest/middlewares"
)

// ICookieManager retrieves the user id from cookies
type ICookieManager interface {
	GetUserID(r *http.Request) string
}

type Handler struct {
	Service *service.Service
	Router  *chi.Mux
	Conf    *config.Config
	Cm      ICookieManager
}

// NewHandler constructor gets a transport instance
func NewHandler(service *service.Service, conf *config.Config) *Handler {
	handler := &Handler{
		Service: service,
		Conf:    conf,
		Cm:      model.CookieManager{},
	}
	handler.Router = handler.InitRoutes()
	return handler
}

// InitRoutes initialization routes
func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middlewares2.GzipHandler)
	router.Use(middlewares2.UserCookieHandler)

	//main routes
	router.Group(func(router chi.Router) {
		router.Get("/{id}", h.HandleGet)
		router.Post("/", h.HandlePost)
		router.Post("/api/shorten", h.HandlePostJSON)
		router.Get("/api/user/urls", h.HandleGetUserURLs)
		router.Get("/ping", h.HandlePing)
		router.Post("/api/shorten/batch", h.HandlePostShortenBatch)
		router.Delete("/api/user/urls", h.HandleDeleteURLs)
	})

	//trusted network routes
	router.Group(func(router chi.Router) {
		trustedNetwork := middlewares2.NewTrustedNetwork(h.Conf)
		router.Use(trustedNetwork.TrustedNetworkHandler)
		router.Get("/api/internal/stats", h.HandleGetStats)
	})

	//Registration of pprof handlers
	router.HandleFunc("/debug/pprof/*", pprof.Index)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))

	return router
}

// HandlePost accepts the URL string in the request body and returns its abbreviated version. Content-Type: text/html
func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
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

	if errors.As(err, &model.ConflictURLError{}) {
		e := err.(model.ConflictURLError)
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

// HandlePostJSON accepts the URL string in the request body and returns its shortened version. Content-Type: application/json
func (h *Handler) HandlePostJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid content type", http.StatusBadRequest)
		return
	}

	originURLReceiver := &struct {
		OriginURL string `json:"url"`
	}{}

	var body []byte
	r.Body.Read(body)

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

	if errors.As(err, &model.ConflictURLError{}) {
		e := err.(model.ConflictURLError)
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

// HandleGet gets the original URL from a shortened link
func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "id")

	url, err := h.Service.GetURLModelByID(shortURL)

	if errors.As(err, &model.GoneError{}) {
		w.WriteHeader(http.StatusGone)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url.Origin)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// HandleGetUserURLs gets all shortened links by the user
func (h *Handler) HandleGetUserURLs(w http.ResponseWriter, r *http.Request) {

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

// HandlePing checks the availability of the database
func (h *Handler) HandlePing(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	err := h.Service.Ping(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandlePostShortenBatch accepts in the request body a set of URLs for abbreviation
func (h *Handler) HandlePostShortenBatch(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputDataList := make([]model.InputCorrelationURL, 0)
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

// HandleDeleteURLs accepts a list of shortened URL to delete
func (h *Handler) HandleDeleteURLs(w http.ResponseWriter, r *http.Request) {
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

// HandleGetStats returns the number of shortened urls and the number of users in the service
func (h *Handler) HandleGetStats(w http.ResponseWriter, _ *http.Request) {
	const nameFunc = "HandleGetStats"
	ctx := context.Background()

	stats, err := h.Service.GetStats(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	statsJSON, err := json.Marshal(stats)
	if err != nil {
		log.Printf("%s error: %s", nameFunc, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(statsJSON)
	if err != nil {
		log.Printf("%s error: %s", nameFunc, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
