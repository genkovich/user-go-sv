package token

type Token struct {
	jwtToken string
}

func NewToken(jwtToken string) *Token {
	return &Token{jwtToken: jwtToken}
}

func (t *Token) GetJwtToken() string {
	return t.jwtToken
}
