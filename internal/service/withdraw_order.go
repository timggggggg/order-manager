package service

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	logpipeline "gitlab.ozon.dev/timofey15g/homework/internal/log_pipeline"
	"gitlab.ozon.dev/timofey15g/homework/internal/metrics"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

func (s *Service) WithdrawOrder(ctx context.Context, req *pb.TReqWithdrawOrder) (*pb.TStringResp, error) {
	logPipeline := logpipeline.GetLogPipelineInstance()

	timer := prometheus.NewTimer(metrics.RequestDuration.WithLabelValues("withdraw_order", "DELETE"))
	defer timer.ObserveDuration()

	order, err := s.storage.WithdrawOrder(ctx, req.OrderID)
	if err != nil {
		metrics.IncrementErrorCounter(err.Error())
		return nil, err
	}

	logPipeline.LogStatusChange(time.Now(), order.ID, models.StatusAccepted, models.StatusWithdrawed)

	return &pb.TStringResp{
		Msg: fmt.Sprintf("order %d withdrawed!", req.OrderID),
	}, nil
}
