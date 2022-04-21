package test

import (
	"fmt"
	"github.com/kotche/url-shortening-service/internal/app/service"
)

type Mock struct {
	original string
	short    string
	err      bool
}

func NewMock(original, short string, err bool) *Mock {
	return &Mock{
		original: original,
		short:    short,
		err:      err,
	}
}

func (m *Mock) Add(url *service.URL) {
	url.SetShort(m.short)
}

func (m *Mock) GetByID(id string) (*service.URL, error) {

	if m.err {
		return nil, fmt.Errorf("key not found")
	}

	return nil, nil
}
