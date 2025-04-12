package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

type WithdrawOrder struct {
	client pb.OrderServiceClient
}

func NewWithdrawOrder(client pb.OrderServiceClient) *WithdrawOrder {
	return &WithdrawOrder{client}
}

func (cmd *WithdrawOrder) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orderIDstr := r.URL.Query().Get("order_id")
	if orderIDstr == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	orderID, err := strconv.ParseInt(orderIDstr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := &pb.TReqWithdrawOrder{
		OrderID: orderID,
	}

	resp, err := cmd.client.WithdrawOrder(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, resp.Msg)
}
