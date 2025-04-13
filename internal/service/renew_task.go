package service

import (
	"context"
	"fmt"

	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

func (s *Service) RenewTask(ctx context.Context, req *pb.TReqRenewTask) (*pb.TStringResp, error) {
	if err := s.outbox.RenewTask(ctx, req.TaskID); err != nil {
		return nil, err
	}

	return &pb.TStringResp{
		Msg: fmt.Sprintf("task %d renewed!", req.TaskID),
	}, nil
}
