package service

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	logpipeline "gitlab.ozon.dev/timofey15g/homework/internal/log_pipeline"
	"gitlab.ozon.dev/timofey15g/homework/internal/metrics"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/internal/packaging"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

func (s *Service) CreateOrder(ctx context.Context, req *pb.TReqAcceptOrder) (*pb.TStringResp, error) {
	logPipeline := logpipeline.GetLogPipelineInstance()

	timer := prometheus.NewTimer(metrics.RequestDuration.WithLabelValues("create_order", "POST"))
	defer timer.ObserveDuration()

	packagingStrategy, err := packaging.NewPackagingStrategy(req.Package, packaging.PackagingStrategies)
	if err != nil {
		return nil, err
	}

	extraPackagingStrategy, err := packaging.NewPackagingStrategy(req.ExtraPackage, packaging.ExtraPackagingStrategies)
	if err != nil {
		return nil, err
	}

	acceptTime := time.Now()
	money, err := models.NewMoney(req.Cost)
	if err != nil {
		return nil, err
	}

	order := models.NewOrder(req.ID, req.UserID, req.StorageDurationDays, acceptTime,
		req.Weight, money, packagingStrategy.Type(), extraPackagingStrategy.Type())

	packageCost, err := validatePackaging(order, packagingStrategy, extraPackagingStrategy)
	if err != nil {
		return nil, err
	}

	order.Cost.Add(packageCost.Amount)

	err = s.storage.CreateOrder(ctx, order)
	if err != nil {
		metrics.IncrementErrorCounter(err.Error())
		return nil, err
	}

	metrics.IncrementOrdersCreated()
	metrics.AddToRevenue(float64(order.Cost.Amount))

	logPipeline.LogStatusChange(time.Now(), order.ID, models.StatusDefault, models.StatusAccepted)

	return &pb.TStringResp{
		Msg: fmt.Sprintf("order %d accepted!", req.ID),
	}, nil
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
