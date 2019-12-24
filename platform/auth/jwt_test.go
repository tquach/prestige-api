package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"testing"
)

func TestGenerateToken(t *testing.T) {
	gen := NewTokenGenerator("secret")

	claimsPayload := map[string]interface{}{}
	claimsPayload["iss"] = "xyz"
	claimsPayload["sid"] = "234234"
	claimsPayload["uid"] = "2902384"
	claimsPayload["exp"] = time.Now().Add(time.Hour * 72).Unix()

	token, err := gen.GenerateAccessToken(claimsPayload)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	parsed, err := jwt.Parse(token, func(tkn *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if !parsed.Valid {
		t.Log(parsed.Raw)
		t.Fail()
	}
}
