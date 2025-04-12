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

func (s *Service) ReturnOrder(ctx context.Context, req *pb.TReqReturnOrder) (*pb.TStringResp, error) {
	logPipeline := logpipeline.GetLogPipelineInstance()

	timer := prometheus.NewTimer(metrics.RequestDuration.WithLabelValues("return_order", "POST"))
	defer timer.ObserveDuration()

	order, err := s.storage.ReturnOrder(ctx, req.OrderID, req.UserID)
	if err != nil {
		metrics.IncrementErrorCounter(err.Error())
		return nil, err
	}

	logPipeline.LogStatusChange(time.Now(), order.ID, models.StatusIssued, models.StatusReturned)

	return &pb.TStringResp{
		Msg: fmt.Sprintf("order %d returned!", req.OrderID),
	}, nil
}
