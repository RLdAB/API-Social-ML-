package domain

import (
	"errors"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrAlreadyFollowing = errors.New("user is already following this seller")
	ErrSelfFollow       = errors.New("cannot follow yourself")
	ErrNotASeller       = errors.New("target user is not a seller")
	//Adicionar outros depois de orientações do Luiz
)
