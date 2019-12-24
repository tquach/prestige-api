// Package users contains user management services.
//
// This file contains functions specific to HTTP routing and handling.
package users

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/tquach/prestige-api/platform/logger"
	"github.com/tquach/prestige-api/platform/render"
	pg "gopkg.in/pg.v4"
)

// Controller handles request routing and servicing.
type Controller struct {
	DB     *pg.DB
	logger logger.Logger
}

// FindUserBySocialID will query the database for a user based on social account id
func (c *Controller) FindUserBySocialID(w http.ResponseWriter, r *http.Request) {
	network := r.URL.Query().Get("network")
	if network != "facebook" {
		render.JSONError(fmt.Errorf("unsupported social network %q", network), http.StatusBadRequest, w)
		return
	}

	socialID := r.URL.Query().Get("socialID")
	c.logger.Debugf("searching for %s user with id %q", network, socialID)
	user := User{}
	if err := c.DB.Model(&user).Where("social_account_id = ?", socialID).Limit(1).Select(); err != nil {
		render.JSONError(fmt.Errorf("no %s user found with id %q)", network, socialID), http.StatusNotFound, w)
		return
	}

	render.JSON(user, http.StatusOK, w)
}

// Find returns the user with this id
func (c *Controller) Find(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		render.JSONError(err, http.StatusBadRequest, w)
	}

	c.logger.Infof("Retrieving user model with id %d", id)
	user := User{ID: id}
	if err := c.DB.Select(&user); err != nil {
		render.JSONError(err, http.StatusNotFound, w)
		return
	}
	render.JSON(user, http.StatusOK, w)
}

// NewController creates a new controller instance.
func NewController(db *pg.DB) *Controller {
	return &Controller{
		DB:     db,
		logger: logger.New("users"),
	}
}
