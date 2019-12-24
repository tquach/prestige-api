package idp

import (
	"errors"
	"net/http"
	"time"

	"github.com/tquach/prestige-api/platform/logger"
)

// OAuthResponse is a response from the IdP. Use this to verify authentication against the IdP.
type OAuthResponse struct {
	UserID        string    `json:"userID"`
	AccessToken   string    `json:"accessToken"`
	ExpiresIn     time.Time `json:"expiresIn,omitempty"`
	SignedRequest string    `json:"signedRequest"`
}

// DefaultIdentityProviders is a mapping of default identity providers supported out of the box.
var DefaultIdentityProviders map[string]IdentityProvider

// ErrorResponse is a possible response from the IdP
type ErrorResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    int    `json:"code"`
}

// IdentityProviderResponse contains the data from doing an identity check.
type IdentityProviderResponse struct {
	ID          string         `json:"id,omitempty"`
	Name        string         `json:"name,omitempty"`
	Error       *ErrorResponse `json:"error,omitempty"`
	RawResponse []byte         `json:"-"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

// IdentityProvider defines an interface for providers to implement.
type IdentityProvider interface {
	ValidateToken(oauthResponse OAuthResponse) (IdentityProviderResponse, error)
}

// New constructs a new instance of an IdP based on the provider network.
func New(network string) (IdentityProvider, error) {
	idp, ok := DefaultIdentityProviders[network]
	if !ok {
		return nil, errors.New("unsupported social network")
	}

	return idp, nil
}

func init() {
	DefaultIdentityProviders = map[string]IdentityProvider{
		"facebook": &FacebookIdentityProvider{
			client: http.Client{},
			logger: logger.New("facebook-idp"),
		},
		"local": &BasicAuthIdentityProvider{
			logger: logger.New("basic-idp"),
		},
	}
}
