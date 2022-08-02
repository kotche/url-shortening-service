package mock

type Generator struct {
	Short string
}

func (m Generator) MakeShortURL() string {
	return m.Short
}
