package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	userApp "github.com/RLdAB/API-Social-ML/internal/user/application"
)

func TestCreateUser_OK(t *testing.T) {
	userRepo := &mockUserRepo{}
	postRepo := &mockPostRepo{}
	userSvc := userApp.NewUserService(userRepo, postRepo)

	handlers := NewUserHandlers(nil, userSvc, nil)

	body := []byte(`{"name":"Alice","is_seller":false}`)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handlers.CreateUser(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d. body=%s", rr.Code, rr.Body.String())
	}
	// opcional: validar que retornou JSON com id/name/is_seller
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	userRepo := &mockUserRepo{}
	postRepo := &mockPostRepo{}
	userSvc := userApp.NewUserService(userRepo, postRepo)

	handlers := NewUserHandlers(nil, userSvc, nil)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handlers.CreateUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d. body=%s", rr.Code, rr.Body.String())
	}
}

func TestCreateUser_MissingName(t *testing.T) {
	userRepo := &mockUserRepo{}
	postRepo := &mockPostRepo{}
	userSvc := userApp.NewUserService(userRepo, postRepo)

	handlers := NewUserHandlers(nil, userSvc, nil)

	body := []byte(`{"is_seller":true}`)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handlers.CreateUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d. body=%s", rr.Code, rr.Body.String())
	}
}

func TestCreateUser_RepoError_InternalServerError(t *testing.T) {
	userRepo := &mockUserRepo{createErr: errors.New("db down")}
	postRepo := &mockPostRepo{}
	userSvc := userApp.NewUserService(userRepo, postRepo)

	handlers := NewUserHandlers(nil, userSvc, nil)

	body := []byte(`{"name":"Alice","is_seller":true}`)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handlers.CreateUser(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d. body=%s", rr.Code, rr.Body.String())
	}
}

func TestCreateUser_BlankName_ServiceInvalidUser(t *testing.T) {
	userRepo := &mockUserRepo{}
	postRepo := &mockPostRepo{}
	userSvc := userApp.NewUserService(userRepo, postRepo)

	handlers := NewUserHandlers(nil, userSvc, nil)

	body := []byte(`{"name":"   ","is_seller":true}`)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handlers.CreateUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d. body=%s", rr.Code, rr.Body.String())
	}
}
