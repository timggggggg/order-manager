//nolint:gocyclo,gocognit
package postgres

import (
	"context"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type PgFacade struct {
	txManager    TransactionManager
	pgRepository *PgRepository
}

func NewPgFacade(txManager TransactionManager, pgRepository *PgRepository) *PgFacade {
	return &PgFacade{
		txManager:    txManager,
		pgRepository: pgRepository,
	}
}

func (s *PgFacade) GetByUserIDCursorPagination(ctx context.Context, userID int64, limit int64, cursorID int64) (models.OrdersSliceStorage, error) {
	result, err := s.pgRepository.GetByUserIDCursorPagination(ctx, userID, limit, cursorID)
	if err != nil {
		return nil, err
	}
	return FromOrdersDBSliceStorage(result), nil
}

func (s *PgFacade) GetAll(ctx context.Context, limit int64, offset int64) (models.OrdersSliceStorage, error) {
	result, err := s.pgRepository.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return FromOrdersDBSliceStorage(result), nil
}

func (s *PgFacade) GetReturnsLimitOffsetPagination(ctx context.Context, limit int64, offset int64) (models.OrdersSliceStorage, error) {
	result, err := s.pgRepository.GetReturnsLimitOffsetPagination(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return FromOrdersDBSliceStorage(result), nil
}

func (s *PgFacade) CreateOrder(ctx context.Context, order *models.Order) error {
	orderDB := ToDTO(order)

	err := s.txManager.RunSerializable(ctx, func(ctxTx context.Context) error {
		if err := s.pgRepository.CreateOrder(ctxTx, orderDB); err != nil {
			return err
		}
		// if err := s.pgRepository.LogOrderEvent(ctxTx, orderDB, models.StatusAccepted); err != nil {
		// 	return err
		// }
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *PgFacade) ReturnOrder(ctx context.Context, orderID int64, userID int64) (order *models.Order, err error) {
	var orderDB *OrderDB
	err = s.txManager.RunReadCommitted(ctx, func(ctxTx context.Context) error {
		orderDB, err = s.pgRepository.GetByID(ctxTx, orderID)
		if err != nil {
			return err
		}

		if err = validateReturn(orderDB, userID); err != nil {
			return err
		}

		orderDB, err = s.returnOrder(ctxTx, orderID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return FromDTO(orderDB), nil
}

func (s *PgFacade) returnOrder(ctx context.Context, id int64) (*OrderDB, error) {
	order, err := s.pgRepository.ReturnOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	// if err := s.pgRepository.LogOrderEvent(ctx, order, models.StatusReturned); err != nil {
	// 	return nil, err
	// }
	return order, nil
}

func (s *PgFacade) IssueOrders(ctx context.Context, ids []int64) (models.OrdersSliceStorage, error) {
	var orders OrdersDBSliceStorage
	err := s.txManager.RunReadCommitted(ctx, func(ctxTx context.Context) error {
		ordersMap, err := s.pgRepository.GetByIDs(ctxTx, ids)
		if err != nil {
			return err
		}

		if err := validateIssues(ordersMap); err != nil {
			return err
		}
		orders, err = s.issueOrders(ctxTx, ids)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return FromOrdersDBSliceStorage(orders), nil
}

func (s *PgFacade) issueOrders(ctx context.Context, ids []int64) (OrdersDBSliceStorage, error) {
	orders, err := s.pgRepository.IssueOrders(ctx, ids)
	if err != nil {
		return nil, err
	}
	// for _, order := range orders {
	// 	if err := s.pgRepository.LogOrderEvent(ctx, order, models.StatusIssued); err != nil {
	// 		return nil, err
	// 	}
	// }
	return orders, nil
}

func (s *PgFacade) WithdrawOrder(ctx context.Context, id int64) (*models.Order, error) {
	var orderDB *OrderDB

	err := s.txManager.RunReadCommitted(ctx, func(ctxTx context.Context) error {
		orderDB, err := s.pgRepository.GetByID(ctxTx, id)
		if err != nil {
			return err
		}

		if err = validateWithdraw(orderDB); err != nil {
			return err
		}
		err = s.withdrawOrder(ctxTx, orderDB)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return FromDTO(orderDB), nil
}

func (s *PgFacade) withdrawOrder(ctx context.Context, order *OrderDB) error {
	if err := s.pgRepository.Delete(ctx, order.ID); err != nil {
		return err
	}
	// if err := s.pgRepository.LogOrderEvent(ctx, order, models.StatusWithdrawed); err != nil {
	// 	return err
	// }
	return nil
}
