package service

import (
	"context"
	"time"

	logpipeline "gitlab.ozon.dev/timofey15g/homework/internal/log_pipeline"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

func (s *Service) IssueOrder(ctx context.Context, req *pb.TReqIssueOrder) (*pb.TStringResp, error) {
	logPipeline := logpipeline.GetLogPipelineInstance()

	orders, err := s.storage.IssueOrders(ctx, req.Ids)
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		logPipeline.LogStatusChange(time.Now(), order.ID, models.StatusAccepted, models.StatusIssued)
	}

	return &pb.TStringResp{
		Msg: "orders issued!",
	}, nil
}
