package commands

import (
	"strconv"
	"time"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/internal/storage"
)

type ReturnOrder struct {
	strg *storage.Storage
}

func NewReturnOrder(strg *storage.Storage) *ReturnOrder {
	return &ReturnOrder{strg}
}

func (cmd *ReturnOrder) Execute(args []string) error {
	if len(args) != 1 {
		return models.ErrorInvalidNumberOfArgs
	}

	orderID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return err
	}
	if orderID <= 0 {
		return models.ErrorNegativeFlag
	}

	order, err := cmd.strg.GetByID(orderID)
	if err != nil {
		return err
	}

	if !time.Now().After(order.ExpireTime) {
		return models.ErrorOrderNotExpired
	}

	if order.Status == storage.Issued {
		return models.ErrorOrderAlreadyIssued
	}

	cmd.strg.DeleteByID(orderID)

	return nil
}
