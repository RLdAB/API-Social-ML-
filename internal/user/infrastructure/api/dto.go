package api

import (
	userDomain "github.com/RLdAB/API-Social-ML/internal/user/domain"
)

type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	IsSeller bool   `json:"is_seller"`
}

func toUserResponse(u userDomain.User) UserResponse {
	return UserResponse{
		ID: u.ID, 
		Name: u.Name, 
		IsSeller: u.IsSeller}
}
