package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/rielj/go-next-chat/internal/token"
	"github.com/rielj/go-next-chat/internal/util"
)

func addAuthorization(
	t *testing.T,
	req *http.Request,
	tokenMaker token.Maker,
	email string,
	duration time.Duration,
) {
	token, payload, err := tokenMaker.CreateToken(email, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	req.Header.Set("x-auth-token", token)

	cookie := new(http.Cookie)
	cookie.Name = "x-auth-token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(duration)
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	req.AddCookie(cookie)
}

func TestAuthMiddleware(t *testing.T) {
	email := util.RandomEmail()

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, email, time.Minute)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name: "InvalidAuthorization",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				req.Header.Set("Authorization", "Bearer invalid-token")
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, email, -time.Hour)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			e := newTestServer(t, nil)
			authPath := "/auth-test"
			e.router.GET(
				authPath,
				func(c echo.Context) error {
					user := c.Get("user").(*jwt.Token)

					mapClaims := user.Claims.(jwt.MapClaims)
					mapClaimsEmail, ok := mapClaims["email"]
					require.True(t, ok)
					require.Equal(t, email, mapClaimsEmail)

					return c.String(http.StatusOK, "Authenticated")
				},
				e.tokenMaker.Middleware(),
			)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, authPath, nil)
			tc.setupAuth(t, req, e.tokenMaker)
			e.router.ServeHTTP(rec, req)
			tc.checkResponse(t, rec)
		})
	}
}
