package commands

import (
	"fmt"
	"strconv"
	"time"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/internal/packaging"
)

type AcceptStorage interface {
	Add(order *models.Order) error
}

type AcceptOrder struct {
	strg AcceptStorage
}

func NewAcceptOrder(strg AcceptStorage) *AcceptOrder {
	return &AcceptOrder{strg}
}

func (cmd *AcceptOrder) Execute(args []string) error {
	if len(args) < 5 {
		return models.ErrorInvalidNumberOfArgs
	}

	iargs, err := parseInts(args[0], args[1], args[2])
	if err != nil {
		return err
	}

	orderID := iargs[0]
	userID := iargs[1]
	storageDurationDays := iargs[2]

	fargs, err := parseFloat(args[3], args[4])
	if err != nil {
		return err
	}
	weight := fargs[0]
	cost := fargs[1]

	// -p packaging -ep extraPackaging
	optionalArgs, err := ParseArgs(args)
	if err != nil {
		return err
	}

	pack, exists := optionalArgs["p"]
	if !exists {
		pack = "film"
	}

	extraPack := optionalArgs["ep"]

	packagingStrategy, err := packaging.NewPackagingStrategy(pack, packaging.PackagingStrategies)
	if err != nil {
		return fmt.Errorf("error creating packaging strategy: %v", err)
	}

	extraPackagingStrategy, err := packaging.NewPackagingStrategy(extraPack, packaging.ExtraPackagingStrategies)
	if err != nil {
		return fmt.Errorf("error creating extra packaging strategy: %v", err)
	}

	acceptTime := time.Now()
	order := models.NewOrder(orderID, userID, storageDurationDays, acceptTime, weight, cost, packagingStrategy.Type(), extraPackagingStrategy.Type())

	packageCost, err := validatePackaging(order, packagingStrategy, extraPackagingStrategy)
	if err != nil {
		return fmt.Errorf("error accepting order: %w", err)
	}

	order.Cost += packageCost

	err = cmd.strg.Add(order)
	if err != nil {
		return err
	}

	fmt.Printf("Order %d accepted!\n", orderID)

	return nil
}

func validatePackaging(order *models.Order, packagingStrategy packaging.Strategy, extraPackagingStrategy packaging.Strategy) (float64, error) {
	if packagingStrategy.Type() == models.PackagingFilm && extraPackagingStrategy.Type() == models.PackagingFilm {
		return 0, models.ErrorPackagingFilmTwice
	}

	packageCost, err := packagingStrategy.CalculateCost(order.Weight)
	if err != nil {
		return 0, err
	}
	extraPackageCost, err := extraPackagingStrategy.CalculateCost(order.Weight)
	if err != nil {
		return 0, err
	}

	return packageCost + extraPackageCost, nil
}

func parseInts(args ...string) ([]int64, error) {
	result := make([]int64, 0)

	for _, s := range args {
		ch, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		if ch <= 0 {
			return nil, models.ErrorNegativeFlag
		}
		result = append(result, ch)
	}
	return result, nil
}

func parseFloat(args ...string) ([]float64, error) {
	result := make([]float64, 0)

	for _, s := range args {
		ch, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		if ch <= 0 {
			return nil, models.ErrorNegativeFlag
		}
		result = append(result, ch)
	}
	return result, nil
}
