package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

type ReturnOrder struct {
	client pb.OrderServiceClient
}

func NewReturnOrder(client pb.OrderServiceClient) *ReturnOrder {
	return &ReturnOrder{client}
}

func (cmd *ReturnOrder) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orderIDstr := r.URL.Query().Get("order_id")
	if orderIDstr == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	userIDstr := r.URL.Query().Get("user_id")
	if userIDstr == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	orderID, err := strconv.ParseInt(orderIDstr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := &pb.TReqReturnOrder{
		OrderID: orderID,
		UserID:  userID,
	}

	resp, err := cmd.client.ReturnOrder(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, resp.Msg)
}
