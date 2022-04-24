package main

import (
	"log"
	"net/http"

	"github.com/kotche/url-shortening-service/internal/app/handler"
	"github.com/kotche/url-shortening-service/internal/app/storage"
	"github.com/kotche/url-shortening-service/internal/config"
)

func main() {

	var UrlStorage handler.Storage = storage.NewUrls()
	handler := handler.NewHandler(UrlStorage)

	log.Fatal(http.ListenAndServe(config.ServerAddr, handler.GetRouter()))
}
