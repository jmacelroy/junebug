package junebug

import (
    "time"

    jwt "github.com/dgrijalva/jwt-go"
)

// ConfluenceTokenGenerator generates conference tokens for auth'ed users.
type ConfluenceTokenGenerator struct {
    Lifetime time.Duration
    Secret   string
    Issuer   string
    Audience string
}

// CreateJWT generates conference tokens for auth'ed users.
func (g ConfluenceTokenGenerator) CreateJWT(qsh string) (string, error) {
    now := time.Now()
    exp := now.Add(g.Lifetime)
    aud := make([]string, 1)
    aud[0] = g.Audience
    claims := jwt.MapClaims{
        "iss": g.Issuer,
        "iat": now.Unix(),
        "exp": exp.Unix(),
        "nbf": now.Unix(),
        "qsh": qsh,
        "aud": aud,
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    return token.SignedString([]byte(g.Secret))
}
