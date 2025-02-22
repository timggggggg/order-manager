package commands

import (
	"fmt"
	"strconv"
	"time"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/internal/storage"
)

type IssueOrder struct {
	strg *storage.Storage
}

func NewIssueOrder(strg *storage.Storage) *IssueOrder {
	return &IssueOrder{strg}
}

func (cmd *IssueOrder) Execute(args []string) error {
	if len(args) < 3 {
		return models.ErrorInvalidNumberOfArgs
	}

	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return err
	}
	if userID <= 0 {
		return models.ErrorNegativeFlag
	}

	orders := make([]*storage.Order, 0)

	for i := 1; i < len(args)-1; i++ {
		orderID, err := strconv.ParseInt(args[i], 10, 64)
		if err != nil {
			return err
		}
		if orderID <= 0 {
			return models.ErrorNegativeFlag
		}

		order, err := cmd.strg.GetByID(orderID)
		if err != nil {
			fmt.Printf("Order %d not found\n", orderID)
			continue
		}
		if order.UserID != userID {
			fmt.Printf("Order %d does not belong to the user %d\n", orderID, userID)
			continue
		}
		orders = append(orders, order)
	}

	mode, err := strconv.ParseInt(args[len(args)-1], 10, 64)
	if err != nil {
		return err
	}

	switch mode {
	case 0:
		// выдать заказы
		for _, order := range orders {
			timeNow := time.Now()

			if timeNow.After(order.ExpireTime) {
				// order.Status = storage.Expired
				fmt.Printf("Order %d expired\n", order.ID)
				continue
			}

			order.IssueTime = timeNow
			order.Status = storage.Issued
			fmt.Printf("Order %d issued!\n", order.ID)
		}
		return cmd.strg.Save()
	case 1:
		// принять возвраты
		for _, order := range orders {
			timeNow := time.Now()

			if order.IssueTime == storage.DefaultTime {
				fmt.Printf("Order %d was not issued\n", order.ID)
				continue
			}

			returnDeadline := order.IssueTime.Add(storage.MaxReturnTime)

			if timeNow.After(returnDeadline) {
				fmt.Printf("Order %d cannot be returned\n", order.ID)
				continue
			}

			order.Status = storage.Returned
			fmt.Printf("Order %d returned!\n", order.ID)
		}
		return cmd.strg.Save()
	default:
		return models.ErrorInvalidIssueMode
	}
}
