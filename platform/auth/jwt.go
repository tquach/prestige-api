package auth

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// TokenService defines the contract.
type TokenService interface {
	// Generates an access token with the given payload.
	GenerateAccessToken(payload map[string]interface{}) (string, error)
	// Verifies the token passed in from a HTTP request.
	VerifyToken(req *http.Request) (*jwt.Token, error)
}

// JWT is a base implementation of TokenService
type JWT struct {
	SecretKey []byte
}

// GenerateAccessToken generates a JWT-compliant access token.
func (j *JWT) GenerateAccessToken(payload map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = payload
	return token.SignedString(j.SecretKey)
}

// VerifyToken parses the token and returns the JSON content. NOTE this does not make any assertions about validity.
func (j *JWT) VerifyToken(req *http.Request) (*jwt.Token, error) {
	return jwt.ParseFromRequest(req, func(t *jwt.Token) (interface{}, error) {
		return j.SecretKey, nil
	})
}

// NewTokenGenerator creates an instance of generator
func NewTokenGenerator(secretKey string) TokenService {
	return &JWT{
		SecretKey: []byte(secretKey),
	}
}
