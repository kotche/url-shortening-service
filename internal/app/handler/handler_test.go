package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kotche/url-shortening-service/internal/app/config"
	mockHandler "github.com/kotche/url-shortening-service/internal/app/handler/mock"
	"github.com/kotche/url-shortening-service/internal/app/service"
	mockService "github.com/kotche/url-shortening-service/internal/app/service/mock"
	mockStorage "github.com/kotche/url-shortening-service/internal/app/storage/mock"
	"github.com/kotche/url-shortening-service/internal/app/storage/test"
	"github.com/kotche/url-shortening-service/internal/app/usecase"
	"github.com/stretchr/testify/assert"
)

func TestHandlerHandleGet(t *testing.T) {

	conf, _ := config.NewConfig()

	type want struct {
		status   int
		location string
	}

	type fields struct {
		id       string
		URL      *usecase.URL
		endpoint string
		err      error
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "url_exists",
			fields: fields{
				URL:      &usecase.URL{Origin: "www.yandex.ru", Short: "qwertyT"},
				endpoint: "/qwertyT",
				id:       "qwertyT",
			},
			want: want{
				status:   http.StatusTemporaryRedirect,
				location: "www.yandex.ru",
			},
		},
		{
			name: "url_not_exists",
			fields: fields{
				URL:      nil,
				endpoint: "/qwertyT",
				id:       "qwertyT",
				err:      errors.New("key not found"),
			},
			want: want{
				status:   http.StatusBadRequest,
				location: "",
			},
		},
		{
			name: "url_gone",
			fields: fields{
				URL:      nil,
				endpoint: "/qwertyT",
				id:       "qwertyT",
				err:      usecase.GoneError{},
			},
			want: want{
				status:   http.StatusGone,
				location: "",
			},
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			control := gomock.NewController(t)
			defer control.Finish()

			repo := mockStorage.NewMockStorage(control)

			repo.EXPECT().GetByID(tt.fields.id).Return(tt.fields.URL, tt.fields.err).Times(1)

			s := service.NewService(repo)
			s.Gen = nil

			h := NewHandler(s, conf)

			r := httptest.NewRequest(http.MethodGet, tt.fields.endpoint, nil)
			w := httptest.NewRecorder()

			h.Router.ServeHTTP(w, r)

			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.want.status, response.StatusCode)
			assert.Equal(t, tt.want.location, response.Header.Get("Location"))

		})
	}
}

func TestHandlerHandlePost(t *testing.T) {

	conf, _ := config.NewConfig()

	type want struct {
		code int
		body string
	}

	type fields struct {
		userID string
		short  string
		origin string
		URLAdd *usecase.URL
		URLGet *usecase.URL
		err    error
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "new_url",
			fields: fields{
				URLAdd: &usecase.URL{Origin: "www.yandex.ru", Short: "qwertyT"},
				URLGet: nil,
				short:  "qwertyT",
				origin: "www.yandex.ru",
				userID: "123",
			},
			want: want{
				code: http.StatusCreated,
				body: "http://localhost:8080/qwertyT",
			},
		},
		{
			name: "conflict_url",
			fields: fields{
				URLAdd: &usecase.URL{Origin: "www.yandex.ru", Short: "qwertyT"},
				URLGet: nil,
				err:    usecase.ConflictURLError{ShortenURL: "qwertyT"},
				short:  "qwertyT",
				origin: "www.yandex.ru",
				userID: "123",
			},
			want: want{
				code: http.StatusConflict,
				body: "http://localhost:8080/qwertyT",
			},
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			control := gomock.NewController(t)
			defer control.Finish()

			repo := mockStorage.NewMockStorage(control)
			repo.EXPECT().Add(tt.fields.userID, tt.fields.URLAdd).Return(tt.fields.err).Times(1)
			repo.EXPECT().GetByID(tt.fields.short).Return(tt.fields.URLGet, tt.fields.err).Times(1)

			generator := mockService.Generator{Short: tt.fields.short}
			cm := mockHandler.CookieManager{Cookie: tt.fields.userID}

			s := service.NewService(repo)
			s.Gen = generator
			h := NewHandler(s, conf)
			h.Cm = cm

			body := bytes.NewBufferString(tt.fields.origin)

			r := httptest.NewRequest(http.MethodPost, "/", body)
			w := httptest.NewRecorder()

			h.Router.ServeHTTP(w, r)

			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.want.code, response.StatusCode)
			assert.Equal(t, tt.want.body, w.Body.String())
		})
	}
}

func TestHandlerHandlePostEmptyBody(t *testing.T) {

	h := NewHandler(nil, nil)

	body := bytes.NewBufferString("")

	r := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()

	h.Router.ServeHTTP(w, r)

	response := w.Result()
	defer response.Body.Close()

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "", w.Body.String())

}

func TestHandlerHandlePostJSON(t *testing.T) {

	conf, _ := config.NewConfig()

	type want struct {
		code int
		body string
	}

	type fields struct {
		userID      string
		short       string
		origin      string
		URLAdd      *usecase.URL
		URLGet      *usecase.URL
		body        string
		err         error
		contentType string
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "new_url_correct",
			fields: fields{
				body:        `{"url":"www.google.com"}`,
				origin:      "www.google.com",
				short:       "qwertyT",
				URLAdd:      &usecase.URL{Origin: "www.google.com", Short: "qwertyT"},
				URLGet:      nil,
				userID:      "123",
				contentType: "application/json",
			},
			want: want{
				code: http.StatusCreated,
				body: `{"result":"http://localhost:8080/qwertyT"}`,
			},
		},
		{
			name: "conflict_url",
			fields: fields{
				body:        `{"url":"www.google.com"}`,
				origin:      "www.google.com",
				short:       "qwertyT",
				URLAdd:      &usecase.URL{Origin: "www.google.com", Short: "qwertyT"},
				URLGet:      nil,
				err:         usecase.ConflictURLError{ShortenURL: "qwertyT"},
				userID:      "123",
				contentType: "application/json",
			},
			want: want{
				code: http.StatusConflict,
				body: `{"result":"http://localhost:8080/qwertyT"}`,
			},
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			control := gomock.NewController(t)
			defer control.Finish()

			repo := mockStorage.NewMockStorage(control)
			repo.EXPECT().Add(tt.fields.userID, tt.fields.URLAdd).Return(tt.fields.err).Times(1)
			repo.EXPECT().GetByID(tt.fields.short).Return(tt.fields.URLGet, tt.fields.err).Times(1)

			generator := mockService.Generator{Short: tt.fields.short}
			cm := mockHandler.CookieManager{Cookie: tt.fields.userID}

			s := service.NewService(repo)
			s.Gen = generator
			h := NewHandler(s, conf)
			h.Cm = cm

			body := bytes.NewBufferString(tt.fields.body)

			r := httptest.NewRequest(http.MethodPost, "/api/shorten", body)
			r.Header.Set("Content-Type", tt.fields.contentType)
			w := httptest.NewRecorder()

			h.Router.ServeHTTP(w, r)

			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.want.code, response.StatusCode)
			assert.Equal(t, tt.want.body, strings.Trim(w.Body.String(), "\n"))
		})
	}
}

func TestHandlerHandlePostJSONBadRequest(t *testing.T) {

	conf, _ := config.NewConfig()

	type want struct {
		code int
	}

	type fields struct {
		short       string
		origin      string
		body        string
		contentType string
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "empty_origin_url",
			fields: fields{
				body:        `{"url":""}`,
				origin:      "",
				short:       "qwertyT",
				contentType: "application/json",
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "another_content_type",
			fields: fields{
				body:        `{"url":"https://www.google.com"}`,
				origin:      "https://www.google.com",
				short:       "qwertyT",
				contentType: "text/plain",
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "wrong_JSON_in_body",
			fields: fields{
				body:        `{"url:"https://www.google.com"}`,
				origin:      "https://www.google.com",
				short:       "qwertyT",
				contentType: "application/json",
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			control := gomock.NewController(t)
			defer control.Finish()

			h := NewHandler(nil, conf)

			body := bytes.NewBufferString(tt.fields.body)

			r := httptest.NewRequest(http.MethodPost, "/api/shorten", body)
			r.Header.Set("Content-Type", tt.fields.contentType)
			w := httptest.NewRecorder()

			h.Router.ServeHTTP(w, r)

			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.want.code, response.StatusCode)
		})
	}
}

func TestHandlerHandleGetURLs(t *testing.T) {

	conf, _ := config.NewConfig()

	type want struct {
		status int
		urls   []*usecase.URL
	}

	type fields struct {
		userID string
		urls   []*usecase.URL
		err    error
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "correct_get_urls",
			fields: fields{
				userID: "123",
				urls: []*usecase.URL{
					{Short: "111", Origin: "www.yandex.ru"},
					{Short: "222", Origin: "www.google.com"},
				},
			},
			want: want{
				status: http.StatusOK,
				urls: []*usecase.URL{
					{Short: "http://localhost:8080/111", Origin: "www.yandex.ru"},
					{Short: "http://localhost:8080/222", Origin: "www.google.com"},
				},
			},
		},
		{
			name: "db_error",
			fields: fields{
				userID: "123",
				err:    errors.New("bd error"),
				urls: []*usecase.URL{
					{Short: "111", Origin: "www.yandex.ru"},
					{Short: "222", Origin: "www.google.com"},
				},
			},
			want: want{
				status: http.StatusInternalServerError,
				urls:   []*usecase.URL{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			control := gomock.NewController(t)
			defer control.Finish()

			repo := mockStorage.NewMockStorage(control)
			repo.EXPECT().GetUserURLs(tt.fields.userID).Return(tt.fields.urls, tt.fields.err).Times(1)

			s := service.NewService(repo)
			s.Gen = nil

			cm := mockHandler.CookieManager{Cookie: tt.fields.userID}
			h := NewHandler(s, conf)
			h.Cm = cm

			r := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			w := httptest.NewRecorder()

			h.Router.ServeHTTP(w, r)

			response := w.Result()
			defer response.Body.Close()

			body, _ := io.ReadAll(response.Body)

			URLResponse := []*usecase.URL{}
			_ = json.Unmarshal(body, &URLResponse)

			assert.Equal(t, tt.want.status, response.StatusCode)
			assert.EqualValues(t, tt.want.urls, URLResponse)
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
			name: "compress_decompress_ok",
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
			name: "no_compress",
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
			name: "bad_compress_header",
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
			name: "bad_compress_body_request",
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

			generator := mockService.Generator{Short: "qwertyT"}
			cm := mockHandler.CookieManager{Cookie: "123"}

			mock := &test.FakeRepo{Short: "qwertyT"}
			s := service.NewService(mock)
			s.Gen = generator

			h := NewHandler(s, conf)
			h.Cm = cm

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
			h.Router.ServeHTTP(w, r)
			response := w.Result()
			defer response.Body.Close()

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
		})
	}
}

func TestHandlePostShortenBatch(t *testing.T) {

	conf, _ := config.NewConfig()
	fakeRepo := &test.FakeRepo{}

	s := service.NewService(fakeRepo)
	s.SetDB(fakeRepo)
	h := NewHandler(s, conf)

	input := []usecase.InputCorrelationURL{
		{
			CorrelationID: "1",
			Origin:        "www.1.ru",
		},
		{
			CorrelationID: "2",
			Origin:        "www.2.ru",
		},
		{
			CorrelationID: "3",
			Origin:        "www.3.ru",
		},
	}

	inputJSON, _ := json.Marshal(input)
	body := bytes.NewBuffer(inputJSON)

	r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", body)
	w := httptest.NewRecorder()

	h.Router.ServeHTTP(w, r)

	response := w.Result()
	defer response.Body.Close()

	assert.Equal(t, http.StatusCreated, response.StatusCode)
}

func TestHandleDeleteURLs(t *testing.T) {

	conf, _ := config.NewConfig()
	fakeRepo := &test.FakeRepo{}

	s := service.NewService(fakeRepo)
	s.SetDB(fakeRepo)
	h := NewHandler(s, conf)

	input := []string{"1", "2", "3"}
	inputJSON, _ := json.Marshal(input)
	body := bytes.NewBuffer(inputJSON)

	r := httptest.NewRequest(http.MethodDelete, "/api/user/urls", body)
	w := httptest.NewRecorder()

	h.Router.ServeHTTP(w, r)

	response := w.Result()
	defer response.Body.Close()

	assert.Equal(t, http.StatusAccepted, response.StatusCode)
}
