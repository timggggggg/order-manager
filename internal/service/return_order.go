package service

import (
	"context"
	"fmt"
	"time"

	logpipeline "gitlab.ozon.dev/timofey15g/homework/internal/log_pipeline"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

func (s *Service) ReturnOrder(ctx context.Context, req *pb.TReqReturnOrder) (*pb.TStringResp, error) {
	logPipeline := logpipeline.GetLogPipelineInstance()

	order, err := s.storage.ReturnOrder(ctx, req.OrderID, req.UserID)
	if err != nil {
		return nil, err
	}

	logPipeline.LogStatusChange(time.Now(), order.ID, models.StatusIssued, models.StatusReturned)

	return &pb.TStringResp{
		Msg: fmt.Sprintf("order %d returned!", req.OrderID),
	}, nil
}
