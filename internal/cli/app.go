package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gitlab.ozon.dev/timofey15g/homework/internal/commands"
	"gitlab.ozon.dev/timofey15g/homework/internal/storage"
)

type App struct {
	storage *storage.Storage
}

func NewApp(storage *storage.Storage) *App {
	return &App{storage}
}

type Command interface {
	Execute(args []string) error
}

func (app *App) Run() {
	reader := bufio.NewReader(os.Stdin)

	cmds := map[string]Command{
		"accept":       commands.NewAcceptOrder(app.storage),
		"return":       commands.NewReturnOrder(app.storage),
		"issue":        commands.NewIssueOrder(app.storage),
		"list_order":   commands.NewListOrder(app.storage),
		"list_return":  commands.NewListReturn(app.storage),
		"list_history": commands.NewListHistory(app.storage),
		"help":         commands.NewHelpCommand(app.storage),
	}

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Printf("Error while reading command: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		args := strings.Split(input, " ")

		if args[0] == "exit" {
			return
		}

		cmd, exists := cmds[args[0]]

		if !exists {
			fmt.Printf("Unknown command\n")
			continue
		}

		err = cmd.Execute(args[1:])
		if err != nil {
			fmt.Printf("Error while executing command: %v\n", err)
			continue
		}
	}
}
