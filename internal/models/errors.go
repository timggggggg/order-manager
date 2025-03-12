package models

import "errors"

var (
	ErrorOrderNotFound      = errors.New("order not found")
	ErrorOrderAlreadyExists = errors.New("order already exists")

	ErrorOrderNotExpired     = errors.New("order cannot be returned to courier")
	ErrorOrderStorageExpired = errors.New("order storage expired")
	ErrorOrderReturnExpired  = errors.New("order cannot be returned to pvz")

	ErrorOrderAlreadyIssued   = errors.New("order already issued")
	ErrorOrderAlreadyReturned = errors.New("order already returned")
	ErrorOrderNotIssued       = errors.New("order is not issued")
	ErrorOrderNotReturned     = errors.New("order is not returned")
	ErrorOrdersDifferentUsers = errors.New("orders belong to multiple users")
	ErrorOrderInvalidUser     = errors.New("order does not belong to the user")

	ErrorNegativeFlag = errors.New("flag cannot be zero or negative")

	ErrorInvalidNumberOfArgs  = errors.New("invalid number of args")
	ErrorInvalidIssueMode     = errors.New("issue mode can only be 0 or 1")
	ErrorInvalidOptionalArgs  = errors.New("invalid optional args")
	ErrorInvalidAmountOfMoney = errors.New("invalid amount of money")

	ErrorPackagingFilmTwice = errors.New("cannot pack with film twice")
)
