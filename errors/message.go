package errors

import "errors"

var (
	ErrBindJson                   = errors.New("BINDING_JSON")
	ErrCreateError                = errors.New("CREATE")
	ErrNotFoundError              = errors.New("NOT_FOUND")
	ErrExchangeOfferError         = errors.New("EXCHANGE_OFFER")
	ErrExchangeOfferAcceptedError = errors.New("EXCHANGE_OFFER_ACCEPT")
	ErrInvalidTokenError          = errors.New("INVALID_TOKEN")
	ErrExpiredTokenError          = errors.New("EXPIRED_TOKEN")
	ErrCreateTokenError           = errors.New("CREATE_TOKEN")
)
