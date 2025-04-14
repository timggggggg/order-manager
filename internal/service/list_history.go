package service

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.ozon.dev/timofey15g/homework/internal/metrics"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

func (s *Service) ListHistory(ctx context.Context, req *pb.TReqListHistory) (*pb.TListResp, error) {

	timer := prometheus.NewTimer(metrics.RequestDuration.WithLabelValues("list_history", "GET"))
	defer timer.ObserveDuration()

	orders, err := s.storage.GetAll(ctx, req.Limit, req.Offset)
	if err != nil {
		metrics.IncrementErrorCounter(err.Error())
		return nil, err
	}

	return &pb.TListResp{
		Orders: models.OrdersSliceModelToProto(orders),
	}, nil
}
