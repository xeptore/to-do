package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type JWT struct {
	secret []byte
}

func New(secret []byte) JWT {
	return JWT{secret}
}

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
		return nil, fmt.Errorf("jwt: failed to generate token: %v", err)
	}

	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.HS512, j.secret))
	if nil != err {
		return nil, fmt.Errorf("jwt: failed to sign token: %v", err)
	}

	return &GenerateResult{Token: string(signed)}, nil
}

type VerificationResult struct {
	UserID string
}

func (j *JWT) VerifyToken(ctx context.Context, token string) (*VerificationResult, error) {
	parseOptions := []jwt.ParseOption{
		jwt.WithVerify(true),
		jwt.WithValidate(true),
		jwt.WithAudience("users"),
		jwt.WithIssuer("https://github.com/xeptore/to-do"),
		jwt.WithAcceptableSkew(time.Second * 5),
		jwt.WithContext(ctx),
		jwt.WithRequiredClaim(jwt.SubjectKey),
		jwt.WithKey(jwa.HS512, []byte(token)),
	}
	parsedToken, err := jwt.Parse([]byte(token), parseOptions...)
	if nil != err {
		// FIXME: differentiate between parsing/validation error, and other kinds of (internal) errors
		return nil, fmt.Errorf("jwt: parsing and verifying token failed: %v", err)
	}

	return &VerificationResult{UserID: parsedToken.Subject()}, nil
}
