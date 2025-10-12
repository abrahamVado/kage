package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"kage/backend/internal/auth"
	"kage/backend/internal/bidding"
	"kage/backend/internal/contracts"
	"kage/backend/internal/trip"
)

// Server bundles HTTP handlers for the backend APIs.
type Server struct {
	arbiter *bidding.Arbiter
	trips   *trip.Manager
	auth    *auth.Validator
}

// NewServer constructs a Server instance.
func NewServer(arbiter *bidding.Arbiter, trips *trip.Manager, validator *auth.Validator) *Server {
	return &Server{arbiter: arbiter, trips: trips, auth: validator}
}

// RegisterRoutes configures Gin routes for REST endpoints.
func (s *Server) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.POST("/bids/evaluate", func(c *gin.Context) {
		if err := s.requireAuth(c); err != nil {
			return
		}

		var payload struct {
			Request contracts.BidRequest `json:"request"`
			Bids    []contracts.Bid      `json:"bids"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 1.- Evaluate incoming bids against rider constraints and capture the winner.
		winner, ok, err := s.arbiter.RankAndSelect(c.Request.Context(), payload.Request, payload.Bids)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "no bids accepted"})
			return
		}

		// 2.- Return the accepted bid to the caller.
		c.JSON(http.StatusOK, gin.H{"winner": winner})
	})

	router.POST("/trips/:id/state", func(c *gin.Context) {
		if err := s.requireAuth(c); err != nil {
			return
		}
		action := struct {
			Action string `json:"action"`
		}{}
		if err := c.ShouldBindJSON(&action); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := s.handleTripAction(c, c.Param("id"), action.Action); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})

	router.GET("/trips/:id/metrics", func(c *gin.Context) {
		if err := s.requireAuth(c); err != nil {
			return
		}
		metrics, ok := s.trips.MetricsFor(c.Param("id"))
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "trip not found"})
			return
		}
		c.JSON(http.StatusOK, metrics)
	})
}

func (s *Server) requireAuth(c *gin.Context) error {
	if s.auth == nil {
		return nil
	}
	token := c.GetHeader("Authorization")
	if err := s.auth.ValidateToken(token); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return err
	}
	return nil
}

func (s *Server) handleTripAction(c *gin.Context, tripID, action string) error {
	switch action {
	case "start":
		return s.trips.StartTrip(c.Request.Context(), tripID)
	case "pause":
		return s.trips.PauseTrip(c.Request.Context(), tripID)
	case "resume":
		return s.trips.ResumeTrip(c.Request.Context(), tripID)
	case "cancel":
		return s.trips.CancelTrip(c.Request.Context(), tripID)
	case "complete":
		return s.trips.CompleteTrip(c.Request.Context(), tripID)
	default:
		return errors.New("unknown action")
	}
}
