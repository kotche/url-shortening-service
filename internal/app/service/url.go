package service

type URL struct {
	Origin string `json:"origin"`
	Short  string `json:"short"`
}

func NewURL(origin, short string) *URL {
	return &URL{Origin: origin,
		Short: short,
	}

}
