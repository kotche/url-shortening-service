package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kotche/url-shortening-service/internal/app/service"
	"github.com/kotche/url-shortening-service/internal/app/storage"
	"github.com/kotche/url-shortening-service/internal/app/storage/test"
	"github.com/kotche/url-shortening-service/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_handleGet(t *testing.T) {

	type want struct {
		code     int
		location string
		body     string
	}

	type fields struct {
		original string
		short    string
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "url exists",
			fields: fields{
				short:    "qwertyT",
				original: "www.yandex.ru",
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: "www.yandex.ru",
				body:     "",
			},
		},
		{
			name: "url not exists",
			fields: fields{
				short:    "qwertyT",
				original: "",
			},
			want: want{
				code:     http.StatusBadRequest,
				location: "",
				body:     "key not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var URLStorage Storage = storage.NewUrls()

			if tt.fields.original != "" {
				URLStorage.Add(service.NewURL(tt.fields.original, tt.fields.short))
			}

			conf := config.NewConfig()
			h := NewHandler(URLStorage, conf)

			r := httptest.NewRequest(http.MethodGet, "/"+tt.fields.short, nil)
			w := httptest.NewRecorder()

			h.GetRouter().ServeHTTP(w, r)

			response := w.Result()

			assert.Equal(t, tt.want.code, response.StatusCode)
			assert.Equal(t, tt.want.location, response.Header.Get("Location"))

			err := response.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.body, w.Body.String())

		})
	}
}

func TestHandler_handlePost(t *testing.T) {

	type want struct {
		code int
		body string
	}

	type fields struct {
		original string
		short    string
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "new url",
			fields: fields{
				original: "www.yandex.ru",
				short:    "qwertyT",
			},
			want: want{
				code: http.StatusCreated,
				body: "http://localhost:8080/qwertyT",
			},
		},
		{
			name: "empty body",
			fields: fields{
				original: "",
				short:    "qwertyT",
			},
			want: want{
				code: http.StatusBadRequest,
				body: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var URLStorage Storage = test.NewMock(tt.fields.original, tt.fields.short)

			conf := config.NewConfig()
			h := NewHandler(URLStorage, conf)

			body := bytes.NewBufferString(tt.fields.original)

			r := httptest.NewRequest(http.MethodPost, "/", body)
			w := httptest.NewRecorder()

			h.GetRouter().ServeHTTP(w, r)

			response := w.Result()

			assert.Equal(t, tt.want.code, response.StatusCode)
			assert.Equal(t, tt.want.body, w.Body.String())

			err := response.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestHandler_handlePostJSON(t *testing.T) {

	type want struct {
		code int
		body string
	}

	type fields struct {
		body        string
		originURL   string
		shortURL    string
		contentType string
		compareBody bool
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "new url correct",
			fields: fields{
				body:        `{"url":"https://www.google.com"}`,
				originURL:   "https://www.google.com",
				shortURL:    "qwertyT",
				contentType: "application/json",
				compareBody: true,
			},
			want: want{
				code: http.StatusCreated,
				body: `{"result":"http://localhost:8080/qwertyT"}`,
			},
		},
		{
			name: "empty origin url",
			fields: fields{
				body:        `{"url":""}`,
				originURL:   "",
				shortURL:    "qwertyT",
				contentType: "application/json",
				compareBody: true,
			},
			want: want{
				code: http.StatusBadRequest,
				body: "URL is empty",
			},
		},
		{
			name: "another content type",
			fields: fields{
				body:        `{"url":"https://www.google.com"}`,
				originURL:   "https://www.google.com",
				shortURL:    "qwertyT",
				contentType: "text/plain",
				compareBody: true,
			},
			want: want{
				code: http.StatusBadRequest,
				body: "Invalid content type",
			},
		},
		{
			name: "wrong JSON in body",
			fields: fields{
				body:        `{"url:"https://www.google.com"}`,
				originURL:   "https://www.google.com",
				shortURL:    "qwertyT",
				contentType: "application/json",
				compareBody: false,
			},
			want: want{
				code: http.StatusBadRequest,
				body: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var URLStorage Storage = test.NewMock(tt.fields.originURL, tt.fields.shortURL)

			conf := config.NewConfig()
			h := NewHandler(URLStorage, conf)

			body := bytes.NewBufferString(tt.fields.body)

			r := httptest.NewRequest(http.MethodPost, "/api/shorten", body)
			r.Header.Set("Content-Type", tt.fields.contentType)
			w := httptest.NewRecorder()

			h.GetRouter().ServeHTTP(w, r)

			response := w.Result()

			assert.Equal(t, tt.want.code, response.StatusCode)

			if tt.fields.compareBody {
				assert.Equal(t, tt.want.body, strings.Trim(w.Body.String(), "\n"))
			}

			err := response.Body.Close()
			require.NoError(t, err)
		})
	}
}
