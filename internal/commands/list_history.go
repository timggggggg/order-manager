package commands

import (
	"fmt"
	"sort"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/internal/storage"
)

type ListHistory struct {
	strg *storage.Storage
}

func NewListHistory(strg *storage.Storage) *ListHistory {
	return &ListHistory{strg}
}

func (cmd *ListHistory) Execute(args []string) error {
	if len(args) != 0 {
		return models.ErrorInvalidNumberOfArgs
	}

	orders := make([]*storage.Order, 0)

	orders = append(orders, cmd.strg.Orders...)

	// по убыванию времени последнего изменения
	sort.Slice(orders, func(i, j int) bool {
		return (orders[i].LastStatusSwitchTime()).After(orders[j].LastStatusSwitchTime())
	})

	for i, order := range orders {
		fmt.Printf("%d) %s\n", i+1, order.String())
	}

	return nil
}
