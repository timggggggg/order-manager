package storage

import (
	"fmt"
	"time"
)

type OrderStatus string

const (
	Default  OrderStatus = ""
	Accepted OrderStatus = "accepted"
	Expired  OrderStatus = "expired"
	Issued   OrderStatus = "issued"
	Returned OrderStatus = "returned"

	MaxReturnTime = time.Hour * 48
)

var (
	DefaultTime time.Time = time.Unix(0, 0)
)

type Order struct {
	ID         int64       `json:"id"`
	UserID     int64       `json:"user_id"`
	Status     OrderStatus `json:"status"`
	AcceptTime time.Time   `json:"accept_time"`
	ExpireTime time.Time   `json:"expire_time"`
	IssueTime  time.Time   `json:"issue_time"`
}

func NewOrder(ID, userID, storageDurationDays int64) *Order {
	acceptTime := time.Now()

	return &Order{
		ID,
		userID,
		Accepted,
		acceptTime,
		acceptTime.AddDate(0, 0, int(storageDurationDays)),
		DefaultTime,
	}
}

func (o *Order) String() string {
	return fmt.Sprintf(
		"Order(ID=%d, UserID=%d, Status=%s, AcceptTime=%s, ExpireTime=%s, IssueTime=%s)",
		o.ID, o.UserID, o.Status, formatTime(o.AcceptTime), formatTime(o.ExpireTime), formatTime(o.IssueTime),
	)
}

func (o *Order) LastStatusSwitchTime() time.Time {
	maxTime := o.AcceptTime
	if o.Status == Expired && o.ExpireTime.After(maxTime) {
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
