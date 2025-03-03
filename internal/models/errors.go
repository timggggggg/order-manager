package models

import "errors"

var (
	ErrorOrderNotFound      = errors.New("order not found")
	ErrorOrderAlreadyExists = errors.New("order already exists")
	ErrorOrderNotExpired    = errors.New("order cannot be returned to courier")
	ErrorOrderAlreadyIssued = errors.New("order already issued")

	ErrorNegativeFlag = errors.New("flag cannot be zero or negative")

	ErrorInvalidNumberOfArgs  = errors.New("invalid number of args")
	ErrorInvalidIssueMode     = errors.New("issue mode can only be 0 or 1")
	ErrorInvalidOptionalArgs  = errors.New("invalid optional args")
	ErrorInvalidAmountOfMoney = errors.New("invalid amount of money")

	ErrorPackagingFilmTwice = errors.New("cannot pack with film twice")
)
