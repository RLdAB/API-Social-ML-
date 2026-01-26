package domain

import "testing"

func TestCanFollow_OK(t *testing.T) {
	u := &User{ID: 1}
	seller := &User{ID: 2, IsSeller: true}

	if err := u.CanFollow(seller); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCanFollow_SelfFollow(t *testing.T) {
	u := &User{ID: 1, IsSeller: false}
	same := &User{ID: 1, IsSeller: true}

	err := u.CanFollow(same)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestCanFollow_TargetNotSeller(t *testing.T) {
	u := &User{ID: 1}
	notSeller := &User{ID: 2, IsSeller: false}

	err := u.CanFollow(notSeller)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
