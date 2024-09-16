package dto

type Token struct {
	value string
}

func NewToken(token string) *Token {
	return &Token{
		value: token,
	}
}

func (t *Token) String() string {
	return t.value
}
