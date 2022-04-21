package service

type URL struct {
	original string
	short    string
}

func NewURL(original, short string) *URL {
	return &URL{original: original,
		short: short,
	}

}

func (u *URL) GetOriginal() string {
	return u.original
}

func (u *URL) GetShort() string {
	return u.short
}

func (u *URL) SetShort(short string) {
	u.short = short
}
