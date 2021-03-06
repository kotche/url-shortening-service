package main

import (
	"log"
	"net/http"

	"github.com/kotche/url-shortening-service/internal/app/handler"
	"github.com/kotche/url-shortening-service/internal/app/service"
	"github.com/kotche/url-shortening-service/internal/app/storage"
	"github.com/kotche/url-shortening-service/internal/app/storage/postgres"
	"github.com/kotche/url-shortening-service/internal/config"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	var (
		URLStorage service.Storage
		Database   service.Database
		serviceURL *service.Service
	)

	if conf.DBConnect != "" {
		Database, err = postgres.NewDB(conf.DBConnect)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		serviceURL = service.NewService(Database)
		serviceURL.SetDB(Database)

		defer func() {
			err = Database.Close()
			if err != nil {
				log.Println(err.Error())
			}
		}()
	} else if conf.FilePath != "" {
		URLStorage, err = storage.NewFileStorage(conf.FilePath)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		serviceURL = service.NewService(URLStorage)
		defer func() {
			err = URLStorage.Close()
			if err != nil {
				log.Println(err.Error())
			}
		}()
	} else {
		URLStorage = storage.NewUrls()
		serviceURL = service.NewService(URLStorage)
	}

	handlerObj := handler.NewHandler(serviceURL, conf)
	log.Fatal(http.ListenAndServe(conf.ServerAddr, handlerObj.GetRouter()))
}
