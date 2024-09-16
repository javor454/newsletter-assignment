package handler

type DecodeToken interface {
	ParseToken(tokenStr string) (string, error)
}

type DecodeTokenHandler struct {
	tokenService DecodeToken
}

func NewDecodeTokenHandler(ts DecodeToken) *DecodeTokenHandler {
	return &DecodeTokenHandler{tokenService: ts}
}

func (d *DecodeTokenHandler) Handle(token string) (string, error) {
	id, err := d.tokenService.ParseToken(token)
	if err != nil {
		return "", err
	}

	return id, nil
}
