package usecase

// URL contains shortened and original URL
type URL struct {
	Short  string `json:"short_url"`
	Origin string `json:"original_url"`
}

func NewURL(origin, short string) *URL {
	return &URL{Origin: origin,
		Short: short,
	}
}

// InputCorrelationURL contains correlation ID and original URL.
type InputCorrelationURL struct {
	CorrelationID string `json:"correlation_id"`
	Origin        string `json:"original_url"`
}

// OutputCorrelationURL contains correlation ID and shortened URL
type OutputCorrelationURL struct {
	CorrelationID string `json:"correlation_id"`
	Short         string `json:"short_url"`
}

// DeleteUserURLs contains user ID and shortened URL
type DeleteUserURLs struct {
	UserID string
	Short  string
}
