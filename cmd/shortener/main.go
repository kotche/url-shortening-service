package main

import (
	"log"
	"net/http"

	"github.com/kotche/url-shortening-service/internal/app/handler"
	"github.com/kotche/url-shortening-service/internal/app/storage"
	"github.com/kotche/url-shortening-service/internal/config"
)

func main() {

	conf := config.NewConfig()
	var URLStorage handler.Storage = storage.NewUrls()
	handler := handler.NewHandler(URLStorage, conf)

	log.Fatal(http.ListenAndServe(conf.GetServerAddr(), handler.GetRouter()))
}
