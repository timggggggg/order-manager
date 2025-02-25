package commands

import (
	"fmt"
	"sort"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type ListHistoryStorage interface {
	GetAllOrders() []*models.Order
}

type ListHistory struct {
	strg ListHistoryStorage
}

func NewListHistory(strg ListHistoryStorage) *ListHistory {
	return &ListHistory{strg}
}

func (cmd *ListHistory) Execute(args []string) error {
	if len(args) != 0 {
		return models.ErrorInvalidNumberOfArgs
	}

	orders := cmd.strg.GetAllOrders()

	// по убыванию времени последнего изменения
	sort.Slice(orders, func(i, j int) bool {
		return (orders[i].LastStatusSwitchTime()).After(orders[j].LastStatusSwitchTime())
	})

	for i, order := range orders {
		fmt.Printf("%d) %s\n", i+1, order.String())
	}

	return nil
}
