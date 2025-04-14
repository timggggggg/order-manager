package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

type OrderJSON struct {
	ID                  int64   `json:"id"`
	UserID              int64   `json:"user_id"`
	StorageDurationDays int64   `json:"storage_duration"`
	Weight              float64 `json:"weight"`
	Cost                string  `json:"cost"`
	Package             string  `json:"package"`
	ExtraPackage        string  `json:"extra_package,omitempty"`
}

type AcceptOrder struct {
	client pb.OrderServiceClient
}

func NewAcceptOrder(client pb.OrderServiceClient) *AcceptOrder {
	return &AcceptOrder{client}
}

func (cmd *AcceptOrder) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var orderJSON OrderJSON
	if err := json.NewDecoder(r.Body).Decode(&orderJSON); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req := &pb.TReqAcceptOrder{
		ID:                  orderJSON.ID,
		UserID:              orderJSON.UserID,
		StorageDurationDays: orderJSON.StorageDurationDays,
		Weight:              orderJSON.Weight,
		Cost:                orderJSON.Cost,
		Package:             orderJSON.Package,
		ExtraPackage:        orderJSON.ExtraPackage,
	}

	resp, err := cmd.client.CreateOrder(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, resp.Msg)
}
