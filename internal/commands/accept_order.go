package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/internal/packaging"
)

type AcceptStorage interface {
	CreateOrder(ctx context.Context, order *models.Order) error
}

type AcceptOrder struct {
	strg AcceptStorage
}

func NewAcceptOrder(strg AcceptStorage) *AcceptOrder {
	return &AcceptOrder{strg}
}

func (cmd *AcceptOrder) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var orderJSON models.OrderJSON
	if err := json.NewDecoder(r.Body).Decode(&orderJSON); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	packagingStrategy, err := packaging.NewPackagingStrategy(orderJSON.Package, packaging.PackagingStrategies)
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating packaging strategy: %v", err), http.StatusBadRequest)
		return
	}

	extraPackagingStrategy, err := packaging.NewPackagingStrategy(orderJSON.ExtraPackage, packaging.ExtraPackagingStrategies)
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating extra packaging strategy: %v", err), http.StatusBadRequest)
		return
	}

	acceptTime := time.Now()
	money, err := models.NewMoney(orderJSON.Cost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	order := models.NewOrder(orderJSON.ID, orderJSON.UserID, orderJSON.StorageDurationDays, acceptTime,
		orderJSON.Weight, money, packagingStrategy.Type(), extraPackagingStrategy.Type())

	packageCost, err := validatePackaging(order, packagingStrategy, extraPackagingStrategy)
	if err != nil {
		http.Error(w, fmt.Sprintf("error validating order: %v", err), http.StatusBadRequest)
		return
	}

	order.Cost.Add(packageCost.Amount)

	err = cmd.strg.CreateOrder(ctx, order)
	if err != nil {
		http.Error(w, fmt.Sprintf("error accepting order: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func validatePackaging(order *models.Order, packagingStrategy packaging.Strategy, extraPackagingStrategy packaging.Strategy) (*models.Money, error) {
	if packagingStrategy.Type() == models.PackagingFilm && extraPackagingStrategy.Type() == models.PackagingFilm {
		return &models.Money{Amount: 0}, models.ErrorPackagingFilmTwice
	}

	packageCost, err := packagingStrategy.CalculateCost(order.Weight)
	if err != nil {
		return &models.Money{Amount: 0}, err
	}
	extraPackageCost, err := extraPackagingStrategy.CalculateCost(order.Weight)
	if err != nil {
		return &models.Money{Amount: 0}, err
	}

	packageCost.Add(extraPackageCost.Amount)

	return packageCost, nil
}
