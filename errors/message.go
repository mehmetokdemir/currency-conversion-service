package errors

import "errors"

var (
	ErrBindJson           = errors.New("BINDING_JSON")
	ErrCreateError        = errors.New("CREATE")
	ErrNotFoundError      = errors.New("NOT_FOUND")
	ErrDataTypeError      = errors.New("TYPE_ASSERTION_ERROR")
	ErrExchangeOfferError = errors.New("EXCHANGE_OFFER")
	ErrSignJwtError       = errors.New("SIGN_JWT")
	ErrInvalidJwtError    = errors.New("INVALID_JWT")
	ErrInvalidTokenError  = errors.New("INVALID_TOKEN")
	ErrExpiredTokenError  = errors.New("EXPIRED_TOKEN")
	ErrCreateTokenError   = errors.New("CREATE_TOKEN")
	ErrInvalidPassword    = errors.New("INVALID_PASSWORD")
)
