package models

import (
	"fmt"
	"time"
)

type OrderStatus string
type OrdersSliceStorage []*Order
type OrdersMapStorage map[int64]*Order

const (
	StatusDefault    OrderStatus = ""
	StatusAccepted   OrderStatus = "accepted"
	StatusExpired    OrderStatus = "expired"
	StatusIssued     OrderStatus = "issued"
	StatusReturned   OrderStatus = "returned"
	StatusWithdrawed OrderStatus = "withdrawed"

	MaxReturnTime = time.Hour * 48
)

var (
	DefaultTime time.Time = time.Unix(0, 0)
)

type Order struct {
	ID           int64
	UserID       int64
	Status       OrderStatus
	AcceptTime   time.Time
	ExpireTime   time.Time
	IssueTime    time.Time
	Weight       float64
	Cost         *Money
	Package      PackagingType
	ExtraPackage PackagingType
}

func NewOrder(ID, userID, storageDurationDays int64, acceptTime time.Time, weight float64, cost *Money, pack, extraPack PackagingType) *Order {
	return &Order{
		ID,
		userID,
		StatusAccepted,
		acceptTime,
		acceptTime.AddDate(0, 0, int(storageDurationDays)),
		DefaultTime,
		weight,
		cost,
		pack,
		extraPack,
	}
}

func (o *Order) String() string {
	return fmt.Sprintf(
		"\t Order(ID=%d, UserID=%d, Status=%s,\n\t AcceptTime=%s, ExpireTime=%s, IssueTime=%s,\n\t Weight=%f, Cost=%s, Packaging=%s, ExtraPackaging=%s)",
		o.ID, o.UserID, o.Status, formatTime(o.AcceptTime), formatTime(o.ExpireTime), formatTime(o.IssueTime), o.Weight, o.Cost.String(), o.Package, o.ExtraPackage,
	)
}

func (o *Order) LastStatusSwitchTime() time.Time {
	maxTime := o.AcceptTime
	if o.Status == StatusExpired && o.ExpireTime.After(maxTime) {
		maxTime = o.ExpireTime
	}
	if o.IssueTime.After(maxTime) {
		maxTime = o.IssueTime
	}
	return maxTime
}

func formatTime(t time.Time) string {
	if t.Equal(DefaultTime) {
		return "nil"
	}
	return t.Format("Monday, 02 January 2006, 15:04:05")
}
