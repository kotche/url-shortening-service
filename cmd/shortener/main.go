package main

import (
	"log"
	"net/http"

	"github.com/kotche/url-shortening-service/internal/app/handler"
	"github.com/kotche/url-shortening-service/internal/app/storage"
	"github.com/kotche/url-shortening-service/internal/config"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	var URLStorage handler.Storage

	if conf.FilePath == "" {
		URLStorage = storage.NewUrls()
	} else {
		URLStorage, err = storage.NewFileStorage(conf.FilePath)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		defer URLStorage.Close()
	}

	handler := handler.NewHandler(URLStorage, conf)

	log.Fatal(http.ListenAndServe(conf.ServerAddr, handler.GetRouter()))
}
