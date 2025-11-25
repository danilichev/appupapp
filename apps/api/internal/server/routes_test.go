package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"apps/api/internal/handlers"
	"apps/api/internal/services"
)

func TestPingHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)

	// Create mock JWT token and claims for authentication
	claims := &services.JwtClaims{
		UserId: "test-user-id",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := &jwt.Token{
		Valid:  true,
		Claims: claims,
	}

	// Set the user token in the context as expected by the handler
	c.Set("user", token)

	if err := handlers.NewPingHandler().GetPing(c); err != nil {
		t.Errorf("handler() error = %v", err)
		return
	}

	if resp.Code != http.StatusOK {
		t.Errorf("handler() wrong status code = %v", resp.Code)
		return
	}
	expected := map[string]string{"message": "pong"}
	var actual map[string]string

	if err := json.NewDecoder(resp.Body).Decode(&actual); err != nil {
		t.Errorf("handler() error decoding response body: %v", err)
		return
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"handler() wrong response body. expected = %v, actual = %v",
			expected,
			actual,
		)
		return
	}
}
