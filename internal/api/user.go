package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	db "github.com/rielj/go-next-chat/internal/db/sqlc"
)

type createUserRequest struct {
	Email          string `json:"email"           validate:"required,email"`
	HashedPassword string `json:"hashed_password" validate:"required"`
	FirstName      string `json:"first_name"      validate:"required"`
	LastName       string `json:"last_name"       validate:"required"`
}

func (s *Server) createUser(c echo.Context) error {
	var input db.CreateUserParams
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := s.store.CreateUser(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return c.JSON(http.StatusOK, user)
}

type getUserRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (s *Server) getUser(c echo.Context) error {
	var input getUserRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := s.store.GetUser(c.Request().Context(), input.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return c.JSON(http.StatusOK, user)
}

type updateUserRequest struct {
	Email     string `json:"email"      validate:"required,email"`
	FirstName string `json:"first_name" validate:"omitempty"`
	LastName  string `json:"last_name"  validate:"omitempty"`
}

func (s *Server) updateUser(c echo.Context) error {
	var input db.UpdateUserParams
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := s.store.UpdateUser(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return c.JSON(http.StatusOK, user)
}
