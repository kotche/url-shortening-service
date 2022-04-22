package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kotche/url-shortening-service/internal/app/service"
	"github.com/kotche/url-shortening-service/internal/app/storage"
	"github.com/kotche/url-shortening-service/internal/app/test"
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

			var UrlStorage storage.Storage = storage.NewUrls()

			if tt.fields.original != "" {
				UrlStorage.Add(service.NewURL(tt.fields.original, tt.fields.short))
			}

			h := NewHandler(UrlStorage)

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

			var UrlStorage storage.Storage = test.NewMock(tt.fields.original, tt.fields.short)

			h := NewHandler(UrlStorage)

			body := bytes.NewBufferString(tt.fields.original)

			r := httptest.NewRequest(http.MethodPost, "/", body)
			w := httptest.NewRecorder()

			h.GetRouter().ServeHTTP(w, r)

			h.handlePost(w, r)
			response := w.Result()

			assert.Equal(t, tt.want.code, response.StatusCode)
			assert.Equal(t, tt.want.body, w.Body.String())

			err := response.Body.Close()
			require.NoError(t, err)
		})
	}
}
