package service

import (
	"context"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

func (s *Service) ListReturns(ctx context.Context, req *pb.TReqListReturns) (*pb.TListResp, error) {

	orders, err := s.storage.GetReturnsLimitOffsetPagination(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	return &pb.TListResp{
		Orders: models.OrdersSliceModelToProto(orders),
	}, nil
}
