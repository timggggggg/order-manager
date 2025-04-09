package models

import "time"

type StatusChangeLog struct {
	OrderID    int64       `json:"order_id"`
	StatusFrom OrderStatus `json:"status_from"`
	StatusTo   OrderStatus `json:"status_to"`
	Timestamp  time.Time   `json:"ts"`
}

type HttpReqLog struct {
	Timestamp   time.Time `json:"ts"`
	Method      string    `json:"method"`
	URL         string    `json:"url"`
	RequestBody string    `json:"request_body"`
}

type HttpRespLog struct {
	Timestamp  time.Time `json:"ts"`
	StatusCode int64     `json:"status_code"`
	Body       string    `json:"body"`
}
