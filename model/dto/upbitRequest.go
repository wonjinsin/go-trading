package dto

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// UpbitTokenPayload ...
type UpbitTokenPayload struct {
	AccessKey    string `json:"access_key"`
	Nonce        string `json:"nonce"`
	QueryHash    string `json:"query_hash"`
	QueryHashAlg string `json:"query_hash_alg"`
}

// NewUpbitTokenPayload ...
func NewUpbitTokenPayload(accessKey string) *UpbitTokenPayload {
	return &UpbitTokenPayload{
		AccessKey: accessKey,
		Nonce:     uuid.New().String(),
	}
}

// GenerateJWT ...
func (p *UpbitTokenPayload) GenerateJWT(secretKey string) (string, error) {
	secretKeyByte := []byte(secretKey)
	tokenClaim := jwt.MapClaims{
		"access_key":     p.AccessKey,
		"nonce":          p.Nonce,
		"query_hash":     p.QueryHash,
		"query_hash_alg": p.QueryHashAlg,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim).
		SignedString(secretKeyByte)
}
