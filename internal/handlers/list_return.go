package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type ListReturnStorage interface {
	GetReturnsLimitOffsetPagination(ctx context.Context, limit int64, offset int64) (models.OrdersSliceStorage, error)
}

type ListReturn struct {
	strg ListReturnStorage
}

func NewListReturn(strg ListReturnStorage) *ListReturn {
	return &ListReturn{strg}
}

func (cmd *ListReturn) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limitstr := r.URL.Query().Get("limit")
	if limitstr == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	offsetstr := r.URL.Query().Get("offset")
	if offsetstr == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	limit, err := strconv.ParseInt(limitstr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset, err := strconv.ParseInt(offsetstr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orders, err := cmd.strg.GetReturnsLimitOffsetPagination(ctx, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		return
	}
}
