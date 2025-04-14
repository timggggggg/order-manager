package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

type ListHistory struct {
	client pb.OrderServiceClient
}

func NewListHistory(client pb.OrderServiceClient) *ListHistory {
	return &ListHistory{client}
}

func (cmd *ListHistory) Execute(w http.ResponseWriter, r *http.Request) {
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

	req := &pb.TReqListHistory{
		Limit:  limit,
		Offset: offset,
	}

	resp, err := cmd.client.ListHistory(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(models.OrdersSliceProtoToModel(resp.Orders))
	if err != nil {
		return
	}
}
