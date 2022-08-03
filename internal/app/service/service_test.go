package service

import (
	"strconv"
	"testing"

	"github.com/kotche/url-shortening-service/internal/app/storage"
	"github.com/kotche/url-shortening-service/internal/app/usecase"
)

func BenchmarkGetURLModel(b *testing.B) {

	const triesN = 1000

	repo := storage.NewUrls()
	s := NewService(repo)

	urls := make([]usecase.URL, triesN)

	for i := 0; i < triesN; i++ {
		str := strconv.Itoa(i)
		urls[i] = usecase.URL{Short: str, Origin: str}
	}

	b.ResetTimer()

	b.Run("current_storage", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < triesN; i++ {
			url, _ := s.GetURLModel(urls[i].Short, urls[i].Origin)
			_ = url
		}
	})
}
