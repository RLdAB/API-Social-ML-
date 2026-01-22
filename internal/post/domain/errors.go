package domain

import "errors"

var (
	ErrInvalidPromotion     = errors.New("invalid promotion data")
	ErrPostNotFound         = errors.New("post not found")
	ErrOnlySellerCanPublish = errors.New("only sellers can publish products")
	ErrInvalidPayload       = errors.New("invalid payload")
	ErrInvalidDate          = errors.New("date must be in format dd-MM-yyyy")
)
