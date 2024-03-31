package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	db "github.com/rielj/go-next-chat/internal/db/sqlc"
	"github.com/rielj/go-next-chat/internal/util"
)

// Server serves HTTP requests for our chatting service.
type Server struct {
	config util.Config
	store  db.Store
	router *echo.Echo
}

// NewServer creates a new server instance.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
		router: NewEcho(),
	}

	server.setupRoutes()

	return server, nil
}

// Setup Echo
func NewEcho() *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	return e
}

// setupRoutes sets up all the routes for our server.
func (s *Server) setupRoutes() {
	s.router.POST("/user", s.createUser)
	s.router.GET("/user", s.getUser)
	s.router.PUT("/user", s.updateUser)
}

// errorResponse is a generic error response.
func errorResponse(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
	}
}

// Start starts the server.
func (s *Server) Start(address string) error {
	return s.router.Start(address)
}
