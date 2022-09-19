package test

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kotche/url-shortening-service/internal/app/config"
	"github.com/kotche/url-shortening-service/internal/app/service"
	mockService "github.com/kotche/url-shortening-service/internal/app/service/mock"
	"github.com/kotche/url-shortening-service/internal/app/storage/test"
	mockHandler "github.com/kotche/url-shortening-service/internal/app/transport/mock"
	"github.com/kotche/url-shortening-service/internal/app/transport/rest"
	"github.com/kotche/url-shortening-service/internal/app/transport/rest/middlewares"
	"github.com/stretchr/testify/assert"
)

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

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			generator := mockService.Generator{Short: "qwertyT"}
			cm := mockHandler.CookieManager{Cookie: "123"}

			mock := &test.FakeRepo{Short: "qwertyT"}
			s := service.NewService(mock)
			s.Gen = generator

			h := rest.NewHandler(s, conf)
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

func TestTrustedNetworkHandler(t *testing.T) {

	type want struct {
		code int
	}

	type header struct {
		name  string
		value string
	}

	tests := []struct {
		name          string
		trustedSubnet string
		header        header
		want          want
	}{
		{
			name:          "access_allowed_x_real_ip",
			trustedSubnet: "192.168.0.1",
			header:        header{name: "X-Real-IP", value: "192.168.0.1"},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:          "access_allowed_x_forward_for",
			trustedSubnet: "192.168.0.1",
			header:        header{name: "X-Forwarded-For", value: "192.168.0.1"},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:          "access_denied",
			trustedSubnet: "192.168.0.1",
			header:        header{name: "X-Real-IP", value: "192.168.0.2"},
			want: want{
				code: http.StatusForbidden,
			},
		},
		{
			name:          "trusted_subnet_empty",
			trustedSubnet: "",
			header:        header{name: "X-Real-IP", value: "192.168.0.1"},
			want: want{
				code: http.StatusForbidden,
			},
		},
		{
			name:          "no_ip_in_headers",
			trustedSubnet: "192.168.0.1",
			want: want{
				code: http.StatusForbidden,
			},
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			conf := &config.Config{TrustedSubnet: tt.trustedSubnet}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

			req := httptest.NewRequest("GET", "http://testing", nil)
			req.Header.Set(tt.header.name, tt.header.value)

			res := httptest.NewRecorder()
			nextHandler(res, req)

			trustedNetwork := middlewares.NewTrustedNetwork(conf)
			handlerToTest := trustedNetwork.TrustedNetworkHandler(nextHandler)

			handlerToTest.ServeHTTP(res, req)
			response := res.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.want.code, response.StatusCode)
		})
	}
}
