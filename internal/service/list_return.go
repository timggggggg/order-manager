package service

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.ozon.dev/timofey15g/homework/internal/metrics"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

func (s *Service) ListReturns(ctx context.Context, req *pb.TReqListReturns) (*pb.TListResp, error) {

	timer := prometheus.NewTimer(metrics.RequestDuration.WithLabelValues("list_returns", "GET"))
	defer timer.ObserveDuration()

	orders, err := s.storage.GetReturnsLimitOffsetPagination(ctx, req.Limit, req.Offset)
	if err != nil {
		metrics.IncrementErrorCounter(err.Error())
		return nil, err
	}

	return &pb.TListResp{
		Orders: models.OrdersSliceModelToProto(orders),
	}, nil
}
