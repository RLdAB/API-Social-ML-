package domain

import (
	"errors"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrAlreadyFollowing = errors.New("user is already following this seller")
	ErrSelfFollow       = errors.New("cannot follow yourself")
	ErrNotASeller       = errors.New("target user is not a seller")
	ErrInvalidPost      = errors.New("invalid post body: userID and content required")
	ErrNotAseller       = errors.New("user must be a seller to create post")
	ErrInvalidUser      = errors.New("user name must not be empty")
	//Adicionar outros depois de orientações do Luiz
)
