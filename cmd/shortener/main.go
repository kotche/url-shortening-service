package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/kotche/url-shortening-service/internal/app/config"
	"github.com/kotche/url-shortening-service/internal/app/handler"
	"github.com/kotche/url-shortening-service/internal/app/server"
	"github.com/kotche/url-shortening-service/internal/app/service"
	"github.com/kotche/url-shortening-service/internal/app/storage"
	"github.com/kotche/url-shortening-service/internal/app/storage/postgres"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {

	printBuildInfo()

	rand.Seed(time.Now().UnixNano())

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

		serviceURL.RunWorker()

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
	srv := server.NewServer(conf, handlerObj.Router)

	if err = srv.Run(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// example: go run -ldflags "-X main.buildVersion=v1.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" main.go
func printBuildInfo() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
