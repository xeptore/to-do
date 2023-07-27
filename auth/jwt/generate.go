package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type GenerateResult struct {
	Token string
}

func (j *JWT) GenerateToken(ctx context.Context, userID string) (*GenerateResult, error) {
	tok, err := jwt.NewBuilder().
		Issuer("https://github.com/xeptore/to-do").
		IssuedAt(time.Now()).
		Audience([]string{"users"}).
		NotBefore(time.Now().Add(3 * time.Second)).
		Build()
	if nil != err {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.HS512, j.secret))
	if nil != err {
		return nil, fmt.Errorf("failed to sign token: %v", err)
	}

	return &GenerateResult{Token: string(signed)}, nil
}
