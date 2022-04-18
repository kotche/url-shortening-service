package main

import (
	"github.com/kotche/url-shortening-service/internal/app/handler"
	"github.com/kotche/url-shortening-service/internal/app/storage"
	"github.com/kotche/url-shortening-service/internal/config"
	"log"
	"net/http"
)

func main() {

	var urls storage.Storage = storage.NewUrls()
	handler := handler.NewHandler(urls)

	mux := http.NewServeMux()
	mux.Handle("/", handler.Handlers())

	log.Fatal(http.ListenAndServe(config.ServerAddr, mux))
}
