package api

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	db "github.com/rielj/go-next-chat/internal/db/sqlc"
	"github.com/rielj/go-next-chat/internal/token"
	"github.com/rielj/go-next-chat/internal/util"
)

// Server serves HTTP requests for our chatting service.
type Server struct {
	config     util.Config
	store      db.Store
	router     *echo.Echo
	tokenMaker token.Maker
}

// NewServer creates a new server instance.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		router:     NewEcho(),
		tokenMaker: tokenMaker,
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
	router := s.router
	router.GET("/user/:email", s.getUser)
	router.POST("/user/login", s.loginUser)

	authRoutes := router.Group("")
	// authRoutes.Use(s.tokenMaker.Middleware())
	authRoutes.Use(authMiddleware(s.tokenMaker))

	authRoutes.POST("/user", s.createUser)
	authRoutes.PUT("/user", s.updateUser)
	authRoutes.GET("/auth-test", func(c echo.Context) error {
		return c.String(http.StatusOK, "Authenticated")
	})
	s.router = router
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
