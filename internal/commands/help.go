package commands

import (
	"fmt"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/internal/storage"
)

type HelpCommand struct {
	strg *storage.Storage
}

func NewHelpCommand(strg *storage.Storage) *HelpCommand {
	return &HelpCommand{strg}
}

func (cmd *HelpCommand) Execute(args []string) error {
	if len(args) != 0 {
		return models.ErrorInvalidNumberOfArgs
	}

	fmt.Println("usage:")
	fmt.Println("\t- accept <orderID> <userID> <storageDurationDays>")
	fmt.Println("\t- return <orderID>")
	fmt.Println("\t- issue <userID> <orderID1> <orderID2> ... <mode>=0/1")
	fmt.Println("\t- list_order <userID> [-n <lastOrdersNumber>] [-s <currentOrderStatus>]")
	fmt.Println("\t- list_return [-p <page>] [-c <ordersPerPage>]")
	fmt.Println("\t- list_history")

	return nil
}
