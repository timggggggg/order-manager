package service

import (
	"context"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

type Storage interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	ReturnOrder(ctx context.Context, orderID int64, userID int64) (order *models.Order, err error)
	IssueOrders(ctx context.Context, ids []int64) (models.OrdersSliceStorage, error)
	WithdrawOrder(ctx context.Context, id int64) (*models.Order, error)
	GetByUserIDCursorPagination(ctx context.Context, userID int64, limit int64, cursorID int64) (models.OrdersSliceStorage, error)
	GetReturnsLimitOffsetPagination(ctx context.Context, limit int64, offset int64) (models.OrdersSliceStorage, error)
	GetAll(ctx context.Context, limit int64, offset int64) (models.OrdersSliceStorage, error)
}

type Outbox interface {
	RenewTask(ctx context.Context, ID int64) error
}

type Service struct {
	service.UnimplementedOrderServiceServer
	storage Storage
	outbox  Outbox
}

func NewService(storage Storage, outbox Outbox) *Service {
	return &Service{storage: storage, outbox: outbox}
}
