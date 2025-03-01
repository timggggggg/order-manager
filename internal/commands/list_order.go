package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nsf/termbox-go"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
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

	ordersTemp := make([]*models.Order, 0)
	for _, order := range cmd.strg.GetAllOrders() {
		if order.UserID == userID {
			ordersTemp = append(ordersTemp, order)
		}
	}

	orders := make([]*models.Order, 0)
	for _, order := range ordersTemp {
		if currentOrderStatus == models.StatusDefault || order.Status == currentOrderStatus {
			orders = append(orders, order)
		}
	}

	if int(lastOrdersNumber) < len(orders) {
		orders = orders[len(orders)-int(lastOrdersNumber):]
	}

	// for i, order := range orders {
	// 	fmt.Printf("%d) %s\n", i+1, order.String())
	// }

	page, pageSize := 0, 3

	err = termbox.Init()
	if err != nil {
		return err
	}
	defer func() {
		termbox.Close()
	}()

	printPage := func() {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		termbox.Flush()

		fmt.Println(strings.Repeat("\n", 50))
		fmt.Println(strings.Repeat("-", 200))

		start := page * pageSize
		end := start + pageSize
		if end > len(orders) {
			end = len(orders)
		}

		for i := start; i < end; i++ {
			fmt.Println(orders[i].String())
		}

		fmt.Printf("\nUserID: %d | Page: %d/%d | 'w' - previous page | 's' - next page | Press 'q' to quit.\n", userID, page+1, (len(orders)+pageSize-1)/pageSize)
		fmt.Println(strings.Repeat("-", 200))
	}

	printPage()
	for {
		switch ev := termbox.PollEvent(); ev.Ch {
		case 's':
			if (page+1)*pageSize < len(orders) {
				page++
				printPage()
			}
		case 'w':
			if page > 0 {
				page--
				printPage()
			}
		case 'q':
			return nil
		default:
		}
	}
}
