package storage

import (
	"time"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type OrderDTO struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Status     string    `json:"status"`
	AcceptTime time.Time `json:"accept_time"`
	ExpireTime time.Time `json:"expire_time"`
	IssueTime  time.Time `json:"issue_time,omitempty"`
}

func ToDTO(o *models.Order) *OrderDTO {
	return &OrderDTO{
		ID:         o.ID,
		UserID:     o.UserID,
		Status:     string(o.Status),
		AcceptTime: o.AcceptTime,
		ExpireTime: o.ExpireTime,
		IssueTime:  o.IssueTime,
	}
}

func FromDTO(d *OrderDTO) *models.Order {
	return &models.Order{
		ID:         d.ID,
		UserID:     d.UserID,
		Status:     models.OrderStatus(d.Status),
		AcceptTime: d.AcceptTime,
		ExpireTime: d.ExpireTime,
		IssueTime:  d.IssueTime,
	}
}
