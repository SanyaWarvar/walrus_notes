package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wn/pkg/apperror"
	"wn/pkg/constants"
	"wn/pkg/response"
	"wn/pkg/token"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func buildAuthToken(t *testing.T) string {
	t.Helper()

	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	layoutID := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	claims := token.CustomClaims{
		UserId:       userID,
		Role:         "ADMIN",
		MainLayoutId: layoutID,
		StandardClaims: jwt.StandardClaims{
			Id: uuid.New().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signed, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}
	return signed
}

func TestRequestIdValidationHandler_Valid(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(constants.RequestIdHeader, uuid.New().String())
	c.Request = req

	RequestIdValidationHandler(c)

	if len(c.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", c.Errors)
	}
	if c.IsAborted() {
		t.Fatal("request should not be aborted")
	}
}

func TestRequestIdValidationHandler_Invalid(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(constants.RequestIdHeader, "not-a-uuid")
	c.Request = req

	RequestIdValidationHandler(c)

	if len(c.Errors) == 0 {
		t.Fatal("expected validation error")
	}
	if !c.IsAborted() {
		t.Fatal("request should be aborted")
	}
}

func TestAuthorizationHandler_Valid(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(constants.AuthorizationHeader, "Bearer "+buildAuthToken(t))
	c.Request = req

	AuthorizationHandler()(c)

	if len(c.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", c.Errors)
	}
	if c.IsAborted() {
		t.Fatal("request should not be aborted")
	}

	role := c.Request.Context().Value(constants.UserRoleCtx)
	if role != "ADMIN" {
		t.Errorf("role = %v, want ADMIN", role)
	}
}

func TestAuthorizationHandler_InvalidHeader(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(constants.AuthorizationHeader, "InvalidFormat")
	c.Request = req

	AuthorizationHandler()(c)

	if len(c.Errors) == 0 {
		t.Fatal("expected authorization error")
	}
	if !c.IsAborted() {
		t.Fatal("request should be aborted")
	}
}

func TestHeaderCtxHandler(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	requestID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set(constants.RequestIdHeader, requestID)
	c.Request = req

	HeaderCtxHandler()(c)

	got := c.Request.Context().Value(constants.RequestIdCtx)
	if got != requestID {
		t.Errorf("RequestIdCtx = %v, want %v", got, requestID)
	}
	gotPath := c.Request.Context().Value(constants.ApiNameCtx)
	if gotPath != "/api/test" {
		t.Errorf("ApiNameCtx = %v, want /api/test", gotPath)
	}
}

func TestErrorHandler(t *testing.T) {
	builder := response.NewResponseBuilder(false)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.Use(ErrorHandler(builder))
	r.GET("/", func(c *gin.Context) {
		_ = c.Error(apperror.NewNotFoundError("missing", "not_found"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusGone {
		t.Errorf("status = %d, want %d", w.Code, http.StatusGone)
	}

	var body response.Response
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if body.Meta.Code != "not_found" {
		t.Errorf("Meta.Code = %q, want not_found", body.Meta.Code)
	}
}
