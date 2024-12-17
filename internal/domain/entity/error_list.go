package entity

import "errors"

var (
	ErrEmptyBrandName    = errors.New("brand name cannot be empty")
	ErrEmptyProductName  = errors.New("product name cannot be empty")
	ErrInvalidPrice      = errors.New("price must be greater than 0")
	ErrInvalidQuantity   = errors.New("quantity must be 0 or greater")
	ErrInvalidBrandID    = errors.New("brand ID is required")
	ErrInvalidAmount     = errors.New("amount must be greater than 0")
	ErrInsufficientStock = errors.New("insufficient stock")
)
