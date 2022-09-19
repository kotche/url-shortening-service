package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/kotche/url-shortening-service/internal/app/config"
	"github.com/kotche/url-shortening-service/internal/app/service"
	"github.com/kotche/url-shortening-service/internal/app/storage"
)

func ExampleHandler_HandlePost() {

	conf, _ := config.NewConfig()

	URLStorage := storage.NewUrls()
	s := service.NewService(URLStorage)
	h := NewHandler(s, conf)

	//Input body:
	//https://www.yandex.ru
	//Return:
	//http://localhost:8080/lBzgbai
	//
	//Response statuses:
	//Success 201 - URL shortened
	//Failure 400 - invalid request format
	//Failure 409 - URL has already been shortened by the current user
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("www.yandex.ru"))
	w := httptest.NewRecorder()

	h.Router.ServeHTTP(w, r)

	response := w.Result()
	defer response.Body.Close()
	fmt.Println(response.StatusCode)

	//Output:
	//201
}

func ExampleHandler_HandlePostJSON() {

	conf, _ := config.NewConfig()

	URLStorage := storage.NewUrls()
	s := service.NewService(URLStorage)
	h := NewHandler(s, conf)

	//Input body:
	//{"url":"www.yandex.ru"}
	//Return:
	//{"result":"http://localhost:8080/lBzgbai"}
	//
	//Response statuses:
	//Success 201 - URL shortened
	//Failure 400 - invalid request format
	//Failure 409 - URL has already been shortened by the current user
	r := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(`{"url":"www.yandex.ru"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Router.ServeHTTP(w, r)

	response := w.Result()
	defer response.Body.Close()
	fmt.Println(response.StatusCode)

	//Output:
	//201
}

func ExampleHandler_HandleGet() {

	conf, _ := config.NewConfig()

	URLStorage := storage.NewUrls()
	s := service.NewService(URLStorage)
	h := NewHandler(s, conf)

	//Endpoint:
	//http://localhost:8080/qwerty
	//Return:
	//https://www.yandex.ru
	//
	//Response statuses:
	//Success 307 - redirect to original url
	//Failure 400 - short URL not found/ internal error
	//Failure 410 - entry deleted
	r := httptest.NewRequest(http.MethodGet, "/qwertyT", nil)
	w := httptest.NewRecorder()

	h.Router.ServeHTTP(w, r)

	response := w.Result()
	defer response.Body.Close()
	fmt.Println(response.StatusCode)

	//Output:
	//400
}
