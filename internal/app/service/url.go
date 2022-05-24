package service

type URL struct {
	Short  string `json:"short_url"`
	Origin string `json:"original_url"`
}

func NewURL(origin, short string) *URL {
	return &URL{Origin: origin,
		Short: short,
	}
}
