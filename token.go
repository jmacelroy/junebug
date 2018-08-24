package junebug

import (
    "time"

    jwt "github.com/dgrijalva/jwt-go"
)

// TokenGenerator generates conference tokens for auth'ed users.
type TokenGenerator struct {
    Lifetime time.Duration
    Secret   string
    Issuer   string
    Audience string
}

// CreateJWT generates conference tokens for auth'ed users.
func (g TokenGenerator) CreateJWT(qsh string) (string, error) {
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

type userClaim struct {
    ID          string `json:"id"`
    DisplayName string `json:"name"`
    AvatarURL   string `json:"avatar"`
}

type contextClaim struct {
    User  userClaim `json:"user"`
    Group string    `json:"group"`
}
