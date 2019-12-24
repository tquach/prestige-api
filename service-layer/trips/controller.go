package trips

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tquach/prestige-api/platform/common"
	"github.com/tquach/prestige-api/platform/logger"
	"github.com/tquach/prestige-api/platform/render"
	pg "gopkg.in/pg.v4"
)

// Controller handles all the requests and responses.
type Controller struct {
	DB     *pg.DB
	logger logger.Logger
}

// ListTrips retrieves all the trips for this user.
func (c *Controller) ListTrips(w http.ResponseWriter, req *http.Request) {
	c.logger.Info("retrieving trips")
	userID := req.Context().Value("userID").(int)

	c.logger.Infof("Querying trips for userId=%d", userID)
	userTrips := []Trip{}
	if err := c.DB.Model(&userTrips).Where("owner_id = ?", userID).
		Column("trip.*", "Destinations").
		Column("trip.*", "Destinations.Ideas").Select(); err != nil {
		render.JSONError(err, http.StatusInternalServerError, w)
		return
	}

	c.logger.Debug("found user trips", userTrips)
	render.JSON(userTrips, http.StatusOK, w)
}

// Save will insert a new trip record in the database.
func (c *Controller) Save(w http.ResponseWriter, req *http.Request) {
	c.logger.Info("saving trips")
	userID := req.Context().Value("userID").(int)
	trip := Trip{
		Status: common.Pending,
	}

	if err := json.NewDecoder(req.Body).Decode(&trip); err != nil {
		render.JSONError(err, http.StatusBadRequest, w)
		return
	}

	if trip.OwnerID != 0 && trip.OwnerID != userID {
		c.logger.Error("cannot create trip for someone else yet")
		render.JSONError(errors.New("ownerID does not match with logged in user"), http.StatusForbidden, w)
		return
	}
	trip.OwnerID = userID

	if err := c.DB.Create(&trip); err != nil {
		c.logger.Error("failed to save trip", err)
		render.JSONError(err, http.StatusInternalServerError, w)
		return
	}

	for _, dest := range trip.Destinations {
		dest.TripID = trip.ID
		if err := c.DB.Create(&dest); err != nil {
			c.logger.Error("failed to save destination", err)
			render.JSONError(err, http.StatusInternalServerError, w)
			return
		}

		for _, idea := range dest.Ideas {
			idea.DestinationID = dest.ID
			if err := c.DB.Create(&idea); err != nil {
				c.logger.Error("failed to save idea", err)
				render.JSONError(err, http.StatusInternalServerError, w)
				return
			}
		}
	}

	render.JSON(trip, http.StatusCreated, w)
}

// Find returns the trip with this id
func (c *Controller) Find(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get(":id"))
	if err != nil {
		render.JSONError(err, http.StatusBadRequest, w)
		return
	}
	userID := req.Context().Value("userID").(int)
	c.logger.Infof("retrieving trip for user %d with tripId: %q", userID, id)

	trip, err := c.findTrip(id)
	if err != nil {
		if err == sql.ErrNoRows {
			render.JSONError(errors.New("no trip found"), http.StatusNotFound, w)
		} else {
			render.JSONError(err, http.StatusInternalServerError, w)
		}
		return
	}
	render.JSON(trip, http.StatusOK, w)
}

// SaveDestination creates a new destination. A trip id is required
func (c *Controller) SaveDestination(w http.ResponseWriter, req *http.Request) {
	destination := Destination{}
	if err := json.NewDecoder(req.Body).Decode(&destination); err != nil {
		render.JSONError(err, http.StatusBadRequest, w)
		return
	}

	if destination.TripID < 0 {
		render.JSONError(errors.New("missing tripId from request body"), http.StatusBadRequest, w)
		return
	}

	if _, err := c.findTrip(destination.TripID); err != nil {
		render.JSONError(err, http.StatusBadRequest, w)
		return
	}

	if err := c.DB.Create(&destination); err != nil {
		render.JSONError(err, http.StatusInternalServerError, w)
		return
	}

	c.logger.Debugf("saved destination %+v", destination)
	render.JSON(destination, http.StatusOK, w)
}

// UpdateDestination deletes a destination. A trip id is required.
func (c *Controller) UpdateDestination(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get(":id"))
	if err != nil {
		render.JSONError(err, http.StatusBadRequest, w)
		return
	}
	c.logger.Debugf("Updating destination with id", id)

	destination := Destination{}
	if err := json.NewDecoder(req.Body).Decode(&destination); err != nil {
		render.JSONError(err, http.StatusBadRequest, w)
		return
	}

	destination.ID = id

	if _, err := c.findDestination(id); err != nil {
		if err == sql.ErrNoRows {
			render.JSONError(fmt.Errorf("no destination found for id %d", id), http.StatusNotFound, w)
		} else {
			c.logger.Errorf("retrieval error %s", err)
			render.JSONError(err, http.StatusInternalServerError, w)
		}
		return
	}

	if err := c.DB.Update(&destination); err != nil {
		render.JSONError(err, http.StatusInternalServerError, w)
		return
	}
	render.JSON(destination, http.StatusOK, w)
}

// SaveIdea adds an idea to existing destination.
func (c *Controller) SaveIdea(w http.ResponseWriter, req *http.Request) {
	idea := Idea{}
	if err := json.NewDecoder(req.Body).Decode(&idea); err != nil {
		render.JSONError(err, http.StatusBadRequest, w)
		return
	}

	if idea.DestinationID < 0 {
		render.JSONError(errors.New("missing destinationId from request body"), http.StatusBadRequest, w)
		return
	}

	if _, err := c.findDestination(idea.DestinationID); err != nil {
		render.JSONError(err, http.StatusBadRequest, w)
		return
	}

	if err := c.DB.Create(&idea); err != nil {
		render.JSONError(err, http.StatusInternalServerError, w)
		return
	}

	c.logger.Debugf("saved idea %+v", idea)
	render.JSON(idea, http.StatusOK, w)
}

func (c *Controller) findTrip(id int) (Trip, error) {
	trip := Trip{}
	if err := c.DB.Model(&trip).Where("id = ?", id).
		Column("trip.*", "Destinations").
		Column("trip.*", "Destinations.Ideas").
		Select(); err != nil {
		return Trip{}, err
	}
	return trip, nil
}

func (c *Controller) findDestination(id int) (Destination, error) {
	destination := Destination{}
	if err := c.DB.Model(&destination).Where("id = ?", id).
		Column("destination.*", "Ideas").
		Select(); err != nil {
		return Destination{}, err
	}
	return destination, nil
}

// NewController creates a new instance of a controller.
func NewController(db *pg.DB) *Controller {
	return &Controller{
		DB:     db,
		logger: logger.New("trips"),
	}
}
