package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/internal/packaging"
)

type OrderJSON struct {
	ID                  int64   `json:"id"`
	UserID              int64   `json:"user_id"`
	StorageDurationDays int64   `json:"storage_duration"`
	Weight              float64 `json:"weight"`
	Cost                string  `json:"cost"`
	Package             string  `json:"package"`
	ExtraPackage        string  `json:"extra_package,omitempty"`
}

type AcceptStorage interface {
	CreateOrder(ctx context.Context, order *models.Order) error
}

type ILogPipeline interface {
	LogStatusChange(TS time.Time, ID int64, statusFrom, statusTo models.OrderStatus)
	Shutdown()
}

type AcceptOrder struct {
	strg        AcceptStorage
	logPipeline ILogPipeline
}

func NewAcceptOrder(strg AcceptStorage, logPipeline ILogPipeline) *AcceptOrder {
	return &AcceptOrder{strg, logPipeline}
}

func (cmd *AcceptOrder) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var orderJSON OrderJSON
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
	cmd.logPipeline.LogStatusChange(time.Now(), order.ID, models.StatusDefault, models.StatusAccepted)
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
