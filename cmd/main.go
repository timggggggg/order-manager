package main

import (
	"flag"
	"fmt"

	"gitlab.ozon.dev/timofey15g/homework/internal/cli"
	"gitlab.ozon.dev/timofey15g/homework/internal/storage"
)

func main() {

	filePath := flag.String("filepath", "./cmd/storage.json", "path to storage file")

	flag.Parse()

	storage, err := storage.NewStorage(*filePath)
	if err != nil {
		fmt.Printf("Error creating storage: %v\n", err)
		return
	}

	app := cli.NewApp(storage)
	app.Run()
}
