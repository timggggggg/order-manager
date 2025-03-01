package commands

import (
	"fmt"
	"strconv"
	"time"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type ReturnStorage interface {
	GetByID(id int64) (*models.Order, error)
	DeleteByID(id int64) error
}

type ReturnOrder struct {
	strg ReturnStorage
}

func NewReturnOrder(strg ReturnStorage) *ReturnOrder {
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

	if order.Status == models.StatusIssued {
		return models.ErrorOrderAlreadyIssued
	}

	cmd.strg.DeleteByID(orderID)
	fmt.Printf("Order %d returned to courier!\n", orderID)

	return nil
}
