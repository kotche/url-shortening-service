package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kotche/url-shortening-service/internal/app/config"
	grpcServer "github.com/kotche/url-shortening-service/internal/app/server/grpc"
	restServer "github.com/kotche/url-shortening-service/internal/app/server/rest"
	"github.com/kotche/url-shortening-service/internal/app/service"
	"github.com/kotche/url-shortening-service/internal/app/storage"
	"github.com/kotche/url-shortening-service/internal/app/storage/postgres"
	"github.com/kotche/url-shortening-service/internal/app/transport/grpc"
	"github.com/kotche/url-shortening-service/internal/app/transport/rest"
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

	ctx, cansel := context.WithCancel(context.Background())
	defer cansel()

	handlerRest := rest.NewHandler(serviceURL, conf)
	restSrv := restServer.NewServer(conf, handlerRest.Router)
	handlerGrpc := grpc.NewHandler(serviceURL, conf)
	grpcSrv := grpcServer.NewServer(conf, handlerGrpc)

	//graceful shutdown
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		log.Println("Starting REST server")
		if err = restSrv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("REST server run error: %s", err)
		}

	}()
	go func() {
		log.Println("Starting gRPC server")
		if err = grpcSrv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("gRCP server run error: %s", err)
		}
	}()

	<-termChan
	log.Println("Servers shutdown")
	cansel()
	if err = restSrv.Stop(ctx); err != nil {
		log.Fatalf("REST server shutdown error: %s", err)
	}
	grpcSrv.Stop()
}

// example: go run -ldflags "-X main.buildVersion=v1.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" main.go
func printBuildInfo() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
