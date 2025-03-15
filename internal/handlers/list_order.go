package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type ListOrderStorage interface {
	GetByUserIDCursorPagination(ctx context.Context, userID int64, limit int64, cursorID int64) (models.OrdersSliceStorage, error)
}

type ListOrder struct {
	strg ListOrderStorage
}

func NewListOrder(strg ListOrderStorage) *ListOrder {
	return &ListOrder{strg}
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

	orders, err := cmd.strg.GetByUserIDCursorPagination(ctx, userID, limit, cursorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		return
	}
}
