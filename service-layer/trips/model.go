package trips

import (
	"time"

	"github.com/tquach/prestige-api/platform/common"
	pg "gopkg.in/pg.v4"
)

// IdeaType classifies an idea.
type IdeaType int

// Idea types
const (
	Place IdeaType = 0 << iota
	Suggestion
)

// String returns a string repr of the IdeaType
func (i IdeaType) String() string {
	switch {
	case i == Place:
		return "Place"
	case i == Suggestion:
		return "Suggestion"
	}
	return "Unknown"
}

// Trip holds the values describing a trip.
type Trip struct {
	ID            int            `json:"id"`
	Name          string         `json:"name"`
	TripDate      time.Time      `json:"tripDate"`
	ShortName     string         `json:"shortName"`
	OwnerID       int            `json:"ownerID"`
	Description   string         `json:"description"`
	Status        common.Status  `json:"status"`
	PermalinkSlug string         `json:"permalinkSlug"`
	LastModified  pg.NullTime    `json:"lastModified" sql:",null"`
	DateCreated   pg.NullTime    `json:"dateCreated" sql:",null"`
	Destinations  []*Destination `pg:"," json:"destinations,omitempty"`
}

// Destination holds destination properties
type Destination struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Position     int           `json:"position"`
	TripID       int           `json:"tripId"`
	Status       common.Status `json:"status"`
	ShortName    string        `json:"shortName"`
	PlaceID      string        `json:"placeId"`
	Latitude     float32       `json:"lat"`
	Longitude    float32       `json:"lng"`
	LastModified pg.NullTime   `json:"lastModified" sql:",null"`
	DateCreated  pg.NullTime   `json:"dateCreated" sql:",null"`
	Ideas        []*Idea       `json:"ideas,omitempty"`
}

// Idea represents a suggestion for the trip. An idea can be a suggested activity, location, venue, etc.
type Idea struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	AuthorID      int           `json:"authorId"`
	Description   string        `json:"description"`
	DestinationID int           `json:"destinationId"`
	PlaceID       string        `json:"placeId"`
	ShortName     string        `json:"shortName"`
	Status        common.Status `json:"status"`
	Latitude      float32       `json:"lat"`
	Longitude     float32       `json:"lng"`
	Type          IdeaType      `json:"type"`
	LastModified  pg.NullTime   `json:"lastModified" sql:",null"`
	DateCreated   pg.NullTime   `json:"dateCreated" sql:",null"`
}

// Comment represents a comment on the idea.
type Comment struct {
	ID           int         `json:"id"`
	AuthorID     string      `json:"author_id"`
	Body         string      `json:"body"`
	LastModified pg.NullTime `json:"lastModified" sql:",null"`
	DateCreated  pg.NullTime `json:"dateCreated" sql:",null"`
}
