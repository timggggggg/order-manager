package service

import (
	"context"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

func (s *Service) ListOrders(ctx context.Context, req *pb.TReqListOrders) (*pb.TListResp, error) {

	orders, err := s.storage.GetByUserIDCursorPagination(ctx, req.UserID, req.Limit, req.CursorID)
	if err != nil {
		return nil, err
	}

	return &pb.TListResp{
		Orders: models.OrdersSliceModelToProto(orders),
	}, nil
}
