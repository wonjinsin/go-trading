package model

import (
	"magmar/util"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Token ...
type Token struct {
	AccessToken  string  `json:"accessToken"`
	RefreshToken *string `json:"refreshToken,omitempty"`
}

// NewToken ...
func NewToken(accessToken string) *Token {
	return &Token{AccessToken: accessToken}
}

// TokenClaim ...
type TokenClaim struct {
	jwt.RegisteredClaims
	Type            string   `json:"type"` // Type of token
	AuthorizedParty string   `json:"azp"`
	Email           string   `json:"email"`
	Roles           []string `json:"scope"`
	Scope           []string `json:"roles"`
}

// NewUserClaim ...
func NewUserClaim(projectName string, domain string, user *User) *TokenClaim {
	return &TokenClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    projectName,
			Subject:   user.UID,
			Audience:  []string{util.TokenAudienceAccount},
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 1)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Type:            util.TokenTypeBaerer,
		AuthorizedParty: projectName,
		Email:           user.Email,
		Roles:           []string{"User"},
		Scope:           []string{"Email"},
	}
}

// GenerateToken ...
func (t *TokenClaim) GenerateToken(prvKey []byte) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		return "", err
	}

	return jwt.NewWithClaims(jwt.SigningMethodRS256, t).SignedString(key)
}
