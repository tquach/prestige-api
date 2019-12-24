package idp

import (
	"errors"

	"github.com/tquach/prestige-api/platform/config"
	"github.com/tquach/prestige-api/platform/logger"
)

// Ensure interface compliance
var _ IdentityProvider = &BasicAuthIdentityProvider{}

// BasicAuthIdentityProvider implements a basic authentication scheme for direct access.
type BasicAuthIdentityProvider struct {
	logger logger.Logger
}

// ValidateToken implements basic auth check.
func (b *BasicAuthIdentityProvider) ValidateToken(oauthResponse OAuthResponse) (IdentityProviderResponse, error) {
	if config.Env != "local" {
		return IdentityProviderResponse{}, errors.New("unauthorized")
	}

	return IdentityProviderResponse{
		ID: oauthResponse.UserID,
	}, nil
}
