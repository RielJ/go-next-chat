package token

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *JWTMaker) CreateToken(
	email string,
	duration time.Duration,
) (string, *Payload, error) {
	payload, err := NewPayload(email, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		fmt.Println(reflect.TypeOf(err), err)
		verr, ok := err.(*jwt.ValidationError)
		fmt.Println(verr.Inner, ok, err)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

// func (maker *JWTMaker) Middleware() echo.MiddlewareFunc {
// 	return echojwt.WithConfig(echojwt.Config{
// 		SigningKey:  []byte(maker.secretKey),
// 		TokenLookup: "cookie:x-auth-token,header:x-auth-token",
// 		Skipper: func(c echo.Context) bool {
// 			paths := []string{"/login", "/register", "/health", "/api/login", "/api/register"}
// 			for _, path := range paths {
// 				if c.Path() == path {
// 					return true
// 				}
// 			}
// 			return false
// 		},
// 	})
// }
