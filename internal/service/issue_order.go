package service

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	logpipeline "gitlab.ozon.dev/timofey15g/homework/internal/log_pipeline"
	"gitlab.ozon.dev/timofey15g/homework/internal/metrics"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

func (s *Service) IssueOrder(ctx context.Context, req *pb.TReqIssueOrder) (*pb.TStringResp, error) {
	logPipeline := logpipeline.GetLogPipelineInstance()

	timer := prometheus.NewTimer(metrics.RequestDuration.WithLabelValues("issue_order", "POST"))
	defer timer.ObserveDuration()

	orders, err := s.storage.IssueOrders(ctx, req.Ids)
	if err != nil {
		metrics.IncrementErrorCounter(err.Error())
		return nil, err
	}

	for _, order := range orders {
		logPipeline.LogStatusChange(time.Now(), order.ID, models.StatusAccepted, models.StatusIssued)
	}

	return &pb.TStringResp{
		Msg: "orders issued!",
	}, nil
}
