package users

import (
	pg "gopkg.in/pg.v4"
)

// User is the base user model.
type User struct {
	ID              int         `json:"id"`
	SocialAccountID string      `json:"socialID"`
	Username        string      `json:"username"`
	FirstName       string      `json:"firstName"`
	LastName        string      `json:"lastName"`
	EmailAddress    string      `json:"emailAddress"`
	LastModified    pg.NullTime `json:"lastModified" sql:",null"`
	DateCreated     pg.NullTime `json:"dateCreated" sql:",null"`
}
