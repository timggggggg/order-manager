package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

type ListOrder struct {
	client pb.OrderServiceClient
}

func NewListOrder(client pb.OrderServiceClient) *ListOrder {
	return &ListOrder{client}
}

func (cmd *ListOrder) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIDstr := r.URL.Query().Get("user_id")
	if userIDstr == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	limitstr := r.URL.Query().Get("limit")
	if limitstr == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	cursorIDstr := r.URL.Query().Get("cursor_id")
	if cursorIDstr == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := strconv.ParseInt(limitstr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cursorID, err := strconv.ParseInt(cursorIDstr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := &pb.TReqListOrders{
		UserID:   userID,
		Limit:    limit,
		CursorID: cursorID,
	}

	resp, err := cmd.client.ListOrders(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(models.OrdersSliceProtoToModel(resp.Orders))
	if err != nil {
		return
	}
}
