package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kotche/url-shortening-service/internal/app/config"
	"github.com/kotche/url-shortening-service/internal/app/model"
	"github.com/kotche/url-shortening-service/internal/app/service"
	mockService "github.com/kotche/url-shortening-service/internal/app/service/mock"
	mockStorage "github.com/kotche/url-shortening-service/internal/app/storage/mock"
	"github.com/kotche/url-shortening-service/internal/app/storage/test"
	mockHandler "github.com/kotche/url-shortening-service/internal/app/transport/mock"
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
		URL      *model.URL
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
				URL:      &model.URL{Origin: "www.yandex.ru", Short: "qwertyT"},
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
				err:      model.GoneError{},
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
			ctx := context.Background()

			repo.EXPECT().GetByID(ctx, tt.fields.id).Return(tt.fields.URL, tt.fields.err).Times(1)

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
		URLAdd *model.URL
		URLGet *model.URL
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
				URLAdd: &model.URL{Origin: "www.yandex.ru", Short: "qwertyT"},
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
				URLAdd: &model.URL{Origin: "www.yandex.ru", Short: "qwertyT"},
				URLGet: nil,
				err:    model.ConflictURLError{ShortenURL: "qwertyT"},
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

			ctx := context.Background()

			repo := mockStorage.NewMockStorage(control)
			repo.EXPECT().Add(ctx, tt.fields.userID, tt.fields.URLAdd).Return(tt.fields.err).Times(1)
			repo.EXPECT().GetByID(ctx, tt.fields.short).Return(tt.fields.URLGet, tt.fields.err).Times(1)

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

	conf := new(config.Config)
	h := NewHandler(nil, conf)

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
		URLAdd      *model.URL
		URLGet      *model.URL
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
				URLAdd:      &model.URL{Origin: "www.google.com", Short: "qwertyT"},
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
				URLAdd:      &model.URL{Origin: "www.google.com", Short: "qwertyT"},
				URLGet:      nil,
				err:         model.ConflictURLError{ShortenURL: "qwertyT"},
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

			ctx := context.Background()

			repo := mockStorage.NewMockStorage(control)
			repo.EXPECT().Add(ctx, tt.fields.userID, tt.fields.URLAdd).Return(tt.fields.err).Times(1)
			repo.EXPECT().GetByID(ctx, tt.fields.short).Return(tt.fields.URLGet, tt.fields.err).Times(1)

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
		urls   []*model.URL
	}

	type fields struct {
		userID string
		urls   []*model.URL
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
				urls: []*model.URL{
					{Short: "111", Origin: "www.yandex.ru"},
					{Short: "222", Origin: "www.google.com"},
				},
			},
			want: want{
				status: http.StatusOK,
				urls: []*model.URL{
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
				urls: []*model.URL{
					{Short: "111", Origin: "www.yandex.ru"},
					{Short: "222", Origin: "www.google.com"},
				},
			},
			want: want{
				status: http.StatusInternalServerError,
				urls:   []*model.URL{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			control := gomock.NewController(t)
			defer control.Finish()

			ctx := context.Background()

			repo := mockStorage.NewMockStorage(control)
			repo.EXPECT().GetUserURLs(ctx, tt.fields.userID).Return(tt.fields.urls, tt.fields.err).Times(1)

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

			URLResponse := []*model.URL{}
			_ = json.Unmarshal(body, &URLResponse)

			assert.Equal(t, tt.want.status, response.StatusCode)
			assert.EqualValues(t, tt.want.urls, URLResponse)
		})
	}
}

func TestHandlePostShortenBatch(t *testing.T) {

	conf, _ := config.NewConfig()
	fakeRepo := &test.FakeRepo{}

	s := service.NewService(fakeRepo)
	s.SetDB(fakeRepo)
	h := NewHandler(s, conf)

	input := []model.InputCorrelationURL{
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

func TestHandleGetStats(t *testing.T) {

	cfg, _ := config.NewConfig()
	cfg.TrustedSubnet = "192.168.1.0"

	type want struct {
		status int
		nUsers int
		nURLs  int
	}

	type fields struct {
		nUsers int
		nURLs  int
		err    error
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "correct_get_stats",
			fields: fields{
				nUsers: 5,
				nURLs:  10,
			},
			want: want{
				status: http.StatusOK,
				nUsers: 5,
				nURLs:  10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx, cansel := context.WithCancel(context.Background())

			control := gomock.NewController(t)
			defer control.Finish()

			repo := mockStorage.NewMockDatabase(control)
			repo.EXPECT().GetNumberOfUsers(ctx).Return(tt.fields.nUsers, tt.fields.err).Times(1)
			repo.EXPECT().GetNumberOfURLs(ctx).Return(tt.fields.nURLs, tt.fields.err).Times(1)

			s := service.NewService(repo)
			s.Gen = nil
			s.SetDB(repo)

			h := NewHandler(s, cfg)

			r := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)
			r.Header.Set("X-Real-IP", cfg.TrustedSubnet)
			w := httptest.NewRecorder()

			h.Router.ServeHTTP(w, r)

			response := w.Result()
			defer response.Body.Close()

			body, _ := io.ReadAll(response.Body)

			var stats model.Stats
			_ = json.Unmarshal(body, &stats)

			assert.Equal(t, tt.want.status, response.StatusCode)
			assert.Equal(t, tt.want.nUsers, stats.NumberOfUsers)
			assert.Equal(t, tt.want.nURLs, stats.NumberOfURLs)

			cansel()
		})
	}
}
