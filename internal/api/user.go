package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	db "github.com/rielj/go-next-chat/internal/db/sqlc"
	"github.com/rielj/go-next-chat/internal/util"
)

type userResponse struct {
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

type createUserRequest struct {
	Email     string `json:"email"      validate:"required,email"`
	Password  string `json:"password"   validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"  validate:"required"`
}

func (s *Server) createUser(c echo.Context) error {
	var input createUserRequest
	if err := validateRequest(c, &input); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	hashedPassword, err := util.HashPassword(input.Password)

	arg := db.CreateUserParams{
		Email:          input.Email,
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		HashedPassword: hashedPassword,
	}

	user, err := s.store.CreateUser(c.Request().Context(), arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return c.JSON(http.StatusConflict, errorResponse(err))
		}
		return c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return c.JSON(http.StatusOK, newUserResponse(user))
}

type getUserRequest struct {
	Email string `param:"email" validate:"required,email"`
}

func (s *Server) getUser(c echo.Context) error {
	var input getUserRequest
	if err := validateRequest(c, &input); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := s.store.GetUser(c.Request().Context(), input.Email)
	if err != nil {
		if err == db.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, errorResponse(err))
		}
		return c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return c.JSON(http.StatusOK, newUserResponse(user))
}

type updateUserRequest struct {
	Email     string `json:"email"      validate:"required,email"`
	FirstName string `json:"first_name" validate:"omitempty"`
	LastName  string `json:"last_name"  validate:"omitempty"`
}

func (s *Server) updateUser(c echo.Context) error {
	var input db.UpdateUserParams
	if err := validateRequest(c, &input); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := s.store.UpdateUser(c.Request().Context(), input)
	if err != nil {
		if err == db.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, errorResponse(err))
		}
		return c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return c.JSON(http.StatusOK, newUserResponse(user))
}

type loginUserRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  int64        `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt int64        `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (s *Server) loginUser(c echo.Context) error {
	var input loginUserRequest
	if err := validateRequest(c, &input); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := s.store.GetUser(c.Request().Context(), input.Email)
	if err != nil {
		if err == db.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, errorResponse(err))
		}
		return c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	if err := util.CheckPassword(input.Password, user.HashedPassword); err != nil {
		return c.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.Email,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.Email,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	session, err := s.store.CreateSession(c.Request().Context(), db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Email:        user.Email,
		RefreshToken: refreshToken,
		UserAgent:    c.Request().UserAgent(),
		ClientIp:     c.RealIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	response := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt.Unix(),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt.Unix(),
		User:                  newUserResponse(user),
	}

	return c.JSON(http.StatusOK, response)
}
