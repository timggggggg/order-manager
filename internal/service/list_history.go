package service

import (
	"context"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

func (s *Service) ListHistory(ctx context.Context, req *pb.TReqListHistory) (*pb.TListResp, error) {

	orders, err := s.storage.GetAll(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	return &pb.TListResp{
		Orders: models.OrdersSliceModelToProto(orders),
	}, nil
}
