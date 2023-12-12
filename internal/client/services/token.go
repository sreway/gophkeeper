package services

type TokenProvider interface {
	GetToken() string
	SetToken(value string)
}

type tokenProvider struct {
	token string
}

func (t *tokenProvider) GetToken() string {
	return t.token
}

func (t *tokenProvider) SetToken(value string) {
	t.token = value
}

func NewTokenProvider() *tokenProvider {
	return new(tokenProvider)
}
