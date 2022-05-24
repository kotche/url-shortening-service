package handler

import (
	"bytes"
	"compress/gzip"
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

	conf, _ := config.NewConfig()

	type want struct {
		code     int
		location string
		body     string
	}

	type fields struct {
		userID   string
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
				userID:   "111",
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
				userID:   "111",
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

			URLStorage := storage.NewUrls()

			if tt.fields.original != "" {
				_ = URLStorage.Add(tt.fields.userID, service.NewURL(tt.fields.original, tt.fields.short))
			}

			service := service.NewService(URLStorage)

			h := NewHandler(service, conf)

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

	conf, _ := config.NewConfig()

	type want struct {
		code int
		body string
	}

	type fields struct {
		original string
		short    string
		userID   string
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
				userID:   "111",
			},
			want: want{
				code: http.StatusBadRequest,
				body: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mock := test.NewMock(tt.fields.original, tt.fields.short)
			service := service.NewService(mock)

			h := NewHandler(service, conf)

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

	conf, _ := config.NewConfig()

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
		userID      string
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

			mock := test.NewMock(tt.fields.originURL, tt.fields.shortURL)

			service := service.NewService(mock)
			h := NewHandler(service, conf)

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

func TestGzipHandle(t *testing.T) {

	conf, _ := config.NewConfig()

	type want struct {
		code int
	}

	type header struct {
		name  string
		value string
	}

	type fields struct {
		compressRequest    bool
		decompressResponse bool
		headers            []header
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "compress/decompress ok",
			fields: fields{
				compressRequest:    true,
				decompressResponse: true,
				headers: []header{
					{name: "Content-Encoding", value: "gzip"},
					{name: "Accept-Encoding", value: "gzip"},
				},
			},
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name: "no compress",
			fields: fields{
				compressRequest:    false,
				decompressResponse: false,
				headers:            []header{},
			},
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name: "bad compress header",
			fields: fields{
				compressRequest:    true,
				decompressResponse: false,
				headers: []header{
					{name: "Content-Encoding", value: "test"},
				},
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "bad compress body request",
			fields: fields{
				compressRequest:    false,
				decompressResponse: false,
				headers: []header{
					{name: "Content-Encoding", value: "gzip"},
				},
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mock := test.NewMock("https://www.yandex.com", "qwertyT")
			service := service.NewService(mock)
			h := NewHandler(service, conf)

			data := []byte(`{"url":"https://www.google.com"}`)

			var body bytes.Buffer

			if tt.fields.compressRequest {
				writer, _ := gzip.NewWriterLevel(&body, gzip.BestSpeed)
				writer.Write(data)
				writer.Close()
			} else {
				body = *bytes.NewBuffer(data)
			}

			r := httptest.NewRequest(http.MethodPost, "/api/shorten", &body)
			r.Header.Set("Content-Type", "application/json")

			for _, v := range tt.fields.headers {
				r.Header.Set(v.name, v.value)
			}

			w := httptest.NewRecorder()
			h.GetRouter().ServeHTTP(w, r)
			response := w.Result()

			assert.Equal(t, tt.want.code, response.StatusCode)

			if tt.fields.decompressResponse {
				reader, _ := gzip.NewReader(bytes.NewReader(w.Body.Bytes()))
				defer reader.Close()
				var b bytes.Buffer
				b.ReadFrom(reader)
				bodyResponse := b.Bytes()

				if !strings.Contains(string(bodyResponse), "qwertyT") {
					t.Error("response body does not match")
				}
			}

			err := response.Body.Close()
			require.NoError(t, err)
		})
	}
}
