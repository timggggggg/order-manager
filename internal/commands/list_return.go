package commands

import (
	"fmt"
	"strconv"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/internal/storage"
)

type ListReturn struct {
	strg *storage.Storage
}

func NewListReturn(strg *storage.Storage) *ListReturn {
	return &ListReturn{strg}
}

func (cmd *ListReturn) Execute(args []string) error {

	// -p page -c ordersPerPage
	optionalArgs, err := ParseArgs(args)
	if err != nil {
		return err
	}

	page, ordersPerPage := int64(1), int64(1)

	pageTemp, exists := optionalArgs["p"]
	if exists {
		page, err = strconv.ParseInt(pageTemp, 10, 64)
		if err != nil {
			return err
		}
	}

	ordersPerPageTemp, exists := optionalArgs["c"]
	if exists {
		ordersPerPage, err = strconv.ParseInt(ordersPerPageTemp, 10, 64)
		if err != nil {
			return err
		}
	}

	offset := (page - 1) * ordersPerPage

	if offset > int64(len(cmd.strg.Orders)) {
		return models.ErrorInvalidOptionalArgs
	}

	orders := make([]*storage.Order, 0)
	for _, order := range cmd.strg.Orders {
		if order.Status == storage.Returned {
			orders = append(orders, order)
		}
	}

	if offset >= int64(len(orders)) {
		return nil
	}

	orders = orders[offset:min(int64(len(orders)), offset+ordersPerPage)]
	for i, order := range orders {
		fmt.Printf("%d) %s\n", i+1, order.String())
	}

	return nil
}
