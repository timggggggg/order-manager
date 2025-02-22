package commands

import (
	"fmt"
	"strconv"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/internal/storage"
)

func ParseArgs(args []string) (map[string]string, error) {
	result := make(map[string]string)

	for i := 0; i+1 < len(args); i++ {
		if args[i][0] == '-' {

			_, exists := result[args[i][1:]]
			if exists {
				return nil, models.ErrorInvalidOptionalArgs
			}
			result[args[i][1:]] = args[i+1]
		}
	}

	return result, nil
}

type ListOrder struct {
	strg *storage.Storage
}

func NewListOrder(strg *storage.Storage) *ListOrder {
	return &ListOrder{strg}
}

func (cmd *ListOrder) Execute(args []string) error {
	if len(args) < 1 {
		return models.ErrorInvalidNumberOfArgs
	}

	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return err
	}
	if userID <= 0 {
		return models.ErrorNegativeFlag
	}

	// -n lastOrdersNumber -s currentOrderStatus
	optionalArgs, err := ParseArgs(args)
	if err != nil {
		return err
	}

	lastOrdersNumber, currentOrderStatus := int64(len(cmd.strg.Orders)), storage.Default

	lastOrdersNumberTemp, exists := optionalArgs["n"]
	if exists {
		lastOrdersNumber, err = strconv.ParseInt(lastOrdersNumberTemp, 10, 64)
		if err != nil {
			return models.ErrorInvalidOptionalArgs
		}
	}

	currentOrderStatusTemp, exists := optionalArgs["s"]
	if exists {
		currentOrderStatus = storage.OrderStatus(currentOrderStatusTemp)
	}

	ordersTemp := make([]*storage.Order, 0)
	for _, order := range cmd.strg.Orders {
		if order.UserID == userID {
			ordersTemp = append(ordersTemp, order)
		}
	}

	orders := make([]*storage.Order, 0)
	for _, order := range ordersTemp {
		if currentOrderStatus == storage.Default || order.Status == currentOrderStatus {
			orders = append(orders, order)
		}
	}

	if int(lastOrdersNumber) < len(orders) {
		orders = orders[len(orders)-int(lastOrdersNumber):]
	}

	for i, order := range orders {
		fmt.Printf("%d) %s\n", i+1, order.String())
	}

	return nil
}
