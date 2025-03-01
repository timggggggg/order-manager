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

	orderID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return err
	}
	if orderID <= 0 {
		return models.ErrorNegativeFlag
	}

	userID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return err
	}
	if userID <= 0 {
		return models.ErrorNegativeFlag
	}

	storageDurationDays, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return err
	}
	if storageDurationDays <= 0 {
		return models.ErrorNegativeFlag
	}

	weight, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		return err
	}
	if weight <= 0 {
		return models.ErrorNegativeFlag
	}

	cost, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		return err
	}
	if cost <= 0 {
		return models.ErrorNegativeFlag
	}

	// -p packaging -ep extraPackaging
	optionalArgs, err := ParseArgs(args)
	if err != nil {
		return err
	}

	pack, exists := optionalArgs["p"]
	if !exists {
		pack = "film"
	}

	extraPack, exists := optionalArgs["ep"]
	if !exists {
		extraPack = ""
	}

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
