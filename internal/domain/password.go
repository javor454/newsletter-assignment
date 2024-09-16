package domain

type Password struct {
	value string
}

func NewPassword(value string) (*Password, error) {
	return &Password{value: value}, nil
}

func (p *Password) String() string {
	return p.value
}
