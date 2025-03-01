package commands

import (
	"fmt"
	"strings"

	"github.com/nsf/termbox-go"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

func performPagination(userID int64, orders []*models.Order) error {
	page, pageSize := 0, 3
	err := termbox.Init()
	if err != nil {
		return err
	}

	defer termbox.Close()

	pageCount := (len(orders) + pageSize - 1) / pageSize

	printPage := func() {
		err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		if err != nil {
			fmt.Printf("Error creating termbox: %v", err)
			return
		}
		termbox.Flush()

		fmt.Println(strings.Repeat("\n", 50))
		fmt.Println(strings.Repeat("-", 200))

		start := page * pageSize
		end := min(start+pageSize, len(orders))

		for i := start; i < end; i++ {
			fmt.Println(orders[i].String())
		}

		fmt.Printf("\nUserID: %d | Page: %d/%d | 'w' - previous page | 's' - next page | Press 'q' to quit.\n", userID, page+1, pageCount)
		fmt.Println(strings.Repeat("-", 200))
	}

	printPage()
	for {
		switch ev := termbox.PollEvent(); ev.Ch {
		case 's':
			page = (page + 1) % pageCount
			printPage()
		case 'w':
			page = (page - 1 + pageCount) % pageCount
			printPage()
		case 'q':
			return nil
		default:
		}
	}
}
