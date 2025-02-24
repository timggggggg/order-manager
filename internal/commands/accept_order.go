package commands

import (
	"fmt"
	"strconv"
	"time"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type AcceptStorage interface {
	Add(order *models.Order) error
}

type AcceptOrder struct {
	strg AcceptStorage
}

func NewAcceptOrder(strg AcceptStorage) *AcceptOrder {
	return &AcceptOrder{strg}
}

func (cmd *AcceptOrder) Execute(args []string) error {
	if len(args) != 3 {
		return models.ErrorInvalidNumberOfArgs
	}

	orderID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return err
	}
	if orderID <= 0 {
		return models.ErrorNegativeFlag
	}

	userID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return err
	}
	if userID <= 0 {
		return models.ErrorNegativeFlag
	}

	storageDurationDays, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return err
	}
	if storageDurationDays <= 0 {
		return models.ErrorNegativeFlag
	}

	acceptTime := time.Now()
	order := models.NewOrder(orderID, userID, storageDurationDays, acceptTime)

	err = cmd.strg.Add(order)
	if err != nil {
		return err
	}

	fmt.Printf("Order %d accepted!\n", orderID)

	return nil
}
