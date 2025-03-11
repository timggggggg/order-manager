package postgres

import (
	"database/sql"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type OrderDB struct {
	ID           int64        `db:"id"`
	UserID       int64        `db:"user_id"`
	Status       string       `db:"order_status"`
	AcceptTime   sql.NullTime `db:"accept_time"`
	ExpireTime   sql.NullTime `db:"expire_time"`
	IssueTime    sql.NullTime `db:"issue_time"`
	Weight       float64      `db:"weight"`
	Cost         string       `db:"cost"`
	Package      string       `db:"package"`
	ExtraPackage string       `db:"extra_package"`
}

func ToDTO(o *models.Order) *OrderDB {
	return &OrderDB{
		ID:           o.ID,
		UserID:       o.UserID,
		Status:       string(o.Status),
		AcceptTime:   sql.NullTime{Time: o.AcceptTime, Valid: true},
		ExpireTime:   sql.NullTime{Time: o.ExpireTime, Valid: true},
		IssueTime:    sql.NullTime{Time: o.IssueTime, Valid: true},
		Weight:       o.Weight,
		Cost:         o.Cost.String(),
		Package:      string(o.Package),
		ExtraPackage: string(o.ExtraPackage),
	}
}

func FromDTO(d *OrderDB) *models.Order {
	cost, _ := models.NewMoney(d.Cost)

	return &models.Order{
		ID:           d.ID,
		UserID:       d.UserID,
		Status:       models.OrderStatus(d.Status),
		AcceptTime:   d.AcceptTime.Time,
		ExpireTime:   d.ExpireTime.Time,
		IssueTime:    d.IssueTime.Time,
		Weight:       d.Weight,
		Cost:         cost,
		Package:      models.PackagingType(d.Package),
		ExtraPackage: models.PackagingType(d.ExtraPackage),
	}
}

type OrdersDBSliceStorage []*OrderDB
type OrdersDBMapStorage map[int64]*OrderDB

func OrdersDBStorageSliceToMap(ordersSlice OrdersDBSliceStorage) OrdersDBMapStorage {
	ordersMap := OrdersDBMapStorage{}
	for _, order := range ordersSlice {
		ordersMap[order.ID] = order
	}
	return ordersMap
}

func FromOrdersDBSliceStorage(ordersSlice OrdersDBSliceStorage) models.OrdersSliceStorage {
	result := make(models.OrdersSliceStorage, 0)
	for _, order := range ordersSlice {
		result = append(result, FromDTO(order))
	}

	return result
}

func FromOrdersDBMapStorage(ordersMap OrdersDBMapStorage) models.OrdersMapStorage {
	result := models.OrdersMapStorage{}
	for _, order := range ordersMap {
		result[order.ID] = FromDTO(order)
	}

	return result
}
