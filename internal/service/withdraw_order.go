package service

import (
	"context"
	"fmt"
	"time"

	logpipeline "gitlab.ozon.dev/timofey15g/homework/internal/log_pipeline"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

func (s *Service) WithdrawOrder(ctx context.Context, req *pb.TReqWithdrawOrder) (*pb.TStringResp, error) {
	logPipeline := logpipeline.GetLogPipelineInstance()

	order, err := s.storage.WithdrawOrder(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	logPipeline.LogStatusChange(time.Now(), order.ID, models.StatusAccepted, models.StatusWithdrawed)

	return &pb.TStringResp{
		Msg: fmt.Sprintf("order %d withdrawed!", req.OrderID),
	}, nil
}
