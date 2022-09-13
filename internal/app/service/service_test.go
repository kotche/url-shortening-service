package service

import (
	"context"
	"strconv"
	"testing"

	"github.com/kotche/url-shortening-service/internal/app/model"
	"github.com/kotche/url-shortening-service/internal/app/storage"
	"github.com/kotche/url-shortening-service/internal/app/storage/test"
)

func BenchmarkGetURLModel(b *testing.B) {

	repo := storage.NewUrls()
	s := NewService(repo)

	urls := make([]model.URL, b.N)

	for i := 0; i < b.N; i++ {
		str := strconv.Itoa(i)
		urls[i] = model.URL{Short: str, Origin: str}
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		url, _ := s.GetURLModel(urls[i].Short, urls[i].Origin)
		_ = url
	}
}

func BenchmarkShortenBatch(b *testing.B) {
	const size = 100

	urlsInput := make([]model.InputCorrelationURL, size)

	for i := 0; i < size; i++ {
		str := strconv.Itoa(i)
		urlsInput[i] = model.InputCorrelationURL{CorrelationID: str, Origin: str}
	}

	ctx := context.Background()

	repo := &test.FakeRepo{}
	s := NewService(repo)
	s.SetDB(repo)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		out, _ := s.ShortenBatch(ctx, "123", urlsInput)
		_ = out
	}
}
