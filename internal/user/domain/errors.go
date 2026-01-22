package domain

import (
	"errors"
)

var (
	ErrSelfFollow       = errors.New("cannot follow yourself")
	ErrNotASeller       = errors.New("target user is not a seller")
	ErrInvalidPost      = errors.New("invalid post body: userID and content required")
	ErrNotAseller       = errors.New("user must be a seller to create post")
	ErrInvalidUser      = errors.New("user name cannot be empty")
	ErrInvalidPromotion = errors.New("promotion is invalid")
	ErrInvalidDate      = errors.New("date invalid")
	//Adicionar outros depois de orientações do Luiz
)
