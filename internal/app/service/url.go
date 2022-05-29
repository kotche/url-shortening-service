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

type InputCorrelationURL struct {
	CorrelationID string `json:"correlation_id"`
	Origin        string `json:"original_url"`
}

type OutputCorrelationURL struct {
	CorrelationID string `json:"correlation_id"`
	Short         string `json:"short_url"`
}
