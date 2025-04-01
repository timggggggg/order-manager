//nolint:gocyclo,gocognit
package postgres

import (
	"context"
	"time"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type ICache interface {
	Get(key int64) *models.Order
	Put(key int64, value *models.Order)
}

type PgFacade struct {
	txManager    TransactionManager
	pgRepository *PgRepository
	cache        ICache
	historyCache *OrderHistoryCache
	timeNow      func() time.Time
}

func NewPgFacade(txManager TransactionManager, pgRepository *PgRepository, cache ICache, timeNow func() time.Time) *PgFacade {
	OrderHistoryCache := NewOrderHistoryCache(time.Duration(2)*time.Minute, pgRepository.GetAll, 100, time.Now)

	go OrderHistoryCache.StartBackgroundRefresh(context.Background())
	defer OrderHistoryCache.Stop()

	return &PgFacade{
		txManager:    txManager,
		pgRepository: pgRepository,
		cache:        cache,
		historyCache: OrderHistoryCache,
		timeNow:      timeNow,
	}
}

func (s *PgFacade) GetByUserIDCursorPagination(ctx context.Context, userID int64, limit int64, cursorID int64) (models.OrdersSliceStorage, error) {
	var result models.OrdersSliceStorage

	err := s.txManager.RunReadCommitted(ctx, func(ctxTx context.Context) error {
		resultTemp, err := s.pgRepository.GetByUserIDCursorPagination(ctx, userID, limit, cursorID)
		if err != nil {
			return err
		}

		result = FromOrdersDBSliceStorage(resultTemp)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PgFacade) GetAll(ctx context.Context, limit int64, offset int64) (models.OrdersSliceStorage, error) {
	var result models.OrdersSliceStorage

	t := s.timeNow()
	if limit+offset <= s.historyCache.Size && s.historyCache.LastUpdated.Add(1*time.Minute).After(t) {
		result = make(models.OrdersSliceStorage, 0)
		ordersMP := s.historyCache.GetHistory()
		for i, o := range ordersMP {
			if i >= limit+offset {
				break
			}
			if i >= offset {
				result = append(result, o)
			}
		}

		return result, nil
	}

	err := s.txManager.RunReadCommitted(ctx, func(ctxTx context.Context) error {
		resultTemp, err := s.pgRepository.GetAll(ctx, limit, offset)

		if err != nil {
			return err
		}

		result = FromOrdersDBSliceStorage(resultTemp)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PgFacade) GetReturnsLimitOffsetPagination(ctx context.Context, limit int64, offset int64) (models.OrdersSliceStorage, error) {
	var result models.OrdersSliceStorage

	err := s.txManager.RunReadCommitted(ctx, func(ctxTx context.Context) error {
		resultTemp, err := s.pgRepository.GetReturnsLimitOffsetPagination(ctx, limit, offset)
		if err != nil {
			return err
		}

		result = FromOrdersDBSliceStorage(resultTemp)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PgFacade) CreateOrder(ctx context.Context, order *models.Order) error {
	orderDB := ToDTO(order)
	err := s.txManager.RunSerializable(ctx, func(ctxTx context.Context) error {
		if err := s.pgRepository.CreateOrder(ctxTx, orderDB); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	s.cache.Put(order.ID, order)

	return nil
}

func (s *PgFacade) ReturnOrder(ctx context.Context, orderID int64, userID int64) (*models.Order, error) {
	var orderDBupdated *OrderDB

	order := s.cache.Get(orderID)
	if order != nil {
		orderDB := ToDTO(order)

		err := validateReturn(orderDB, userID)
		if err != nil {
			return nil, err
		}

		err = s.txManager.RunReadCommitted(ctx, func(ctxTx context.Context) error {
			orderDBupdated, err = s.returnOrder(ctxTx, orderID)
			return err
		})

		if err != nil {
			return nil, err
		}

		resultOrder := FromDTO(orderDBupdated)

		s.cache.Put(orderID, resultOrder)

		return resultOrder, nil
	}

	err := s.txManager.RunReadCommitted(ctx, func(ctxTx context.Context) error {
		orderDB, err := s.pgRepository.GetByID(ctxTx, orderID)
		if err != nil {
			return err
		}

		if err = validateReturn(orderDB, userID); err != nil {
			return err
		}

		orderDBupdated, err = s.returnOrder(ctxTx, orderID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	resultOrder := FromDTO(orderDBupdated)
	s.cache.Put(orderID, resultOrder)

	return resultOrder, nil
}

func (s *PgFacade) returnOrder(ctx context.Context, id int64) (*OrderDB, error) {
	var order *OrderDB

	order, err := s.pgRepository.ReturnOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *PgFacade) getOrdersMap(ids []int64) (OrdersDBMapStorage, bool) {
	ordersMap := make(OrdersDBMapStorage)
	for _, id := range ids {
		o := s.cache.Get(id)
		if o == nil {
			return nil, false
		}
		ordersMap[id] = ToDTO(o)
	}

	return ordersMap, true
}

func (s *PgFacade) IssueOrders(ctx context.Context, ids []int64) (models.OrdersSliceStorage, error) {
	var orders OrdersDBSliceStorage

	ordersMap, succsessfulCacheLookup := s.getOrdersMap(ids)

	if succsessfulCacheLookup {
		err := validateIssues(ordersMap)
		if err != nil {
			return nil, err
		}
		err = s.txManager.RunReadCommitted(ctx, func(ctxTx context.Context) error {
			orders, err = s.issueOrders(ctxTx, ids)
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return nil, err
		}

		result := FromOrdersDBSliceStorage(orders)

		for _, o := range result {
			s.cache.Put(o.ID, o)
		}

		return result, nil
	}

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

	result := FromOrdersDBSliceStorage(orders)

	for _, o := range result {
		s.cache.Put(o.ID, o)
	}

	return result, nil
}

func (s *PgFacade) issueOrders(ctx context.Context, ids []int64) (OrdersDBSliceStorage, error) {
	var orders OrdersDBSliceStorage

	orders, err := s.pgRepository.IssueOrders(ctx, ids)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *PgFacade) WithdrawOrder(ctx context.Context, id int64) (*models.Order, error) {
	var order *models.Order

	order = s.cache.Get(id)
	if order != nil {
		orderDB := ToDTO(order)

		err := validateWithdraw(orderDB)
		if err != nil {
			return nil, err
		}

		err = s.txManager.RunReadCommitted(ctx, func(ctxTx context.Context) error {
			err = s.withdrawOrder(ctxTx, orderDB)
			return err
		})

		if err != nil {
			return nil, err
		}

		return order, nil
	}

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

		order = FromDTO(orderDB)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *PgFacade) withdrawOrder(ctx context.Context, order *OrderDB) error {
	if err := s.pgRepository.Delete(ctx, order.ID); err != nil {
		return err
	}

	return nil
}
