package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

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
		return nil, fmt.Errorf("jwt: parsing and verifying token failed: %v", err)
	}

	return &VerificationResult{UserID: parsedToken.Subject()}, nil
}
