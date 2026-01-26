package application

import (
	"testing"

	userDomain "github.com/RLdAB/API-Social-ML/internal/user/domain"
)

func TestUserService_CreateUser_OK(t *testing.T) {
	userRepo := &mockUserRepo{}
	postRepo := &mockPostRepo{}
	svc := NewUserService(userRepo, postRepo)

	u := &userDomain.User{Name: "Alice", IsSeller: false}

	if err := svc.CreateUser(u); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(userRepo.created) != 1 {
		t.Fatalf("expected 1 created user, got %d", len(userRepo.created))
	}
	if userRepo.created[0].Name != "Alice" {
		t.Fatalf("expected name Alice, got %s", userRepo.created[0].Name)
	}
}

func TestUserService_CreateUser_InvalidName_Empty(t *testing.T) {
	userRepo := &mockUserRepo{}
	postRepo := &mockPostRepo{}
	svc := NewUserService(userRepo, postRepo)

	u := &userDomain.User{Name: ""}

	err := svc.CreateUser(u)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != userDomain.ErrInvalidUser {
		t.Fatalf("expected ErrInvalidUser, got %v", err)
	}
	if len(userRepo.created) != 0 {
		t.Fatalf("expected no repo call, but got %d", len(userRepo.created))
	}
}

func TestUserService_CreateUser_InvalidName_Blank(t *testing.T) {
	userRepo := &mockUserRepo{}
	postRepo := &mockPostRepo{}
	svc := NewUserService(userRepo, postRepo)

	u := &userDomain.User{Name: "   "}

	err := svc.CreateUser(u)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != userDomain.ErrInvalidUser {
		t.Fatalf("expected ErrInvalidUser, got %v", err)
	}
}
