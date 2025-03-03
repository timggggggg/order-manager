package commands

import (
	"strconv"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type ListOrderStorage interface {
	GetAllOrders() []*models.Order
	GetSize() int64
}

type ListOrder struct {
	strg ListOrderStorage
}

func NewListOrder(strg ListOrderStorage) *ListOrder {
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

	lastOrdersNumber, currentOrderStatus := cmd.strg.GetSize(), models.StatusDefault

	lastOrdersNumberTemp, exists := optionalArgs["n"]
	if exists {
		lastOrdersNumber, err = strconv.ParseInt(lastOrdersNumberTemp, 10, 64)
		if err != nil {
			return models.ErrorInvalidOptionalArgs
		}
	}

	currentOrderStatusTemp, exists := optionalArgs["s"]
	if exists {
		currentOrderStatus = models.OrderStatus(currentOrderStatusTemp)
	}

	orders := filterOrders(cmd.strg.GetAllOrders(), userID, currentOrderStatus)

	if int(lastOrdersNumber) < len(orders) {
		orders = orders[len(orders)-int(lastOrdersNumber):]
	}

	err = performPagination(userID, orders)
	if err != nil {
		return err
	}

	return nil
}

func filterOrders(orders []*models.Order, userID int64, currentOrderStatus models.OrderStatus) []*models.Order {
	ordersTemp := make([]*models.Order, 0)
	for _, order := range orders {
		if userID == 0 || order.UserID == userID {
			ordersTemp = append(ordersTemp, order)
		}
	}

	result := make([]*models.Order, 0)
	for _, order := range ordersTemp {
		if currentOrderStatus == models.StatusDefault || order.Status == currentOrderStatus {
			orders = append(orders, order)
		}
	}

	return result
}
