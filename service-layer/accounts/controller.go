package accounts

import (
	"encoding/json"
	"net/http"
	"time"

	pg "gopkg.in/pg.v4"

	"github.com/tquach/prestige-api/platform/auth"
	"github.com/tquach/prestige-api/platform/idp"
	"github.com/tquach/prestige-api/platform/logger"
	"github.com/tquach/prestige-api/platform/render"
)

// Controller contains all the attributes needed for a full controller.
type Controller struct {
	DB           *pg.DB
	logger       logger.Logger
	TokenService auth.TokenService
}

// RegisterAccount creates a new account for the new user.
func (c *Controller) RegisterAccount(w http.ResponseWriter, r *http.Request) {
	acct := Account{}
	if err := json.NewDecoder(r.Body).Decode(&acct); err != nil {
		render.JSONError(err, http.StatusBadRequest, w)
		return
	}

	c.logger.Infof("Creating new user %v", acct)

	provider, err := idp.New(acct.SocialNetwork)
	if err != nil {
		c.logger.Errorf("invalid social network: %s\n", err.Error())
		render.JSONError(err, http.StatusBadRequest, w)
		return
	}

	idpResponse, err := provider.ValidateToken(acct.OAuthResponse)
	if err != nil {
		c.logger.Error("invalid oauth response")
		render.JSONError(err, http.StatusInternalServerError, w)
		return
	}

	if idpResponse.Error != nil {
		c.logger.Error("error from idp", idpResponse.Error.Message)
		render.JSON(idpResponse.Error, http.StatusUnauthorized, w)
		return
	}

	if err := c.DB.Create(&acct.User); err != nil {
		c.logger.Error(err)
		render.JSONError(err, http.StatusInternalServerError, w)
		return
	}

	claimsPayload := map[string]interface{}{}
	claimsPayload["sid"] = acct.User.SocialAccountID
	claimsPayload["uid"] = acct.User.ID
	claimsPayload["exp"] = time.Now().Add(time.Hour * 72).Unix()

	accessToken, err := c.TokenService.GenerateAccessToken(claimsPayload)
	if err != nil {
		c.logger.Error(err)
		render.JSONError(err, http.StatusInternalServerError, w)
		return
	}

	response := RegistrationResponse{
		User:        acct.User,
		AccessToken: accessToken,
	}
	render.JSON(response, http.StatusOK, w)
}

// NewController creates a new instance of a Controller.
func NewController(db *pg.DB, tokenService auth.TokenService) *Controller {
	return &Controller{
		DB:           db,
		TokenService: tokenService,
		logger:       logger.New("account"),
	}
}
