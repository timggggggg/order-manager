package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

type RenewTask struct {
	client pb.OrderServiceClient
}

func NewRenewTask(client pb.OrderServiceClient) *RenewTask {
	return &RenewTask{client}
}

func (cmd *RenewTask) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id is not provided", http.StatusBadRequest)
		return
	}

	ID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := &pb.TReqRenewTask{TaskID: ID}

	resp, err := cmd.client.RenewTask(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, resp.Msg)
}
