package jwt

type JWT struct {
	secret []byte
}

func New(secret []byte) JWT {
	return JWT{secret}
}
