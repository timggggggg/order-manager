package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

type IssueOrder struct {
	client pb.OrderServiceClient
}

func NewIssueOrder(client pb.OrderServiceClient) *IssueOrder {
	return &IssueOrder{client}
}

func (cmd *IssueOrder) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var ids []int64
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req := &pb.TReqIssueOrder{
		Ids: ids,
	}

	resp, err := cmd.client.IssueOrder(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, resp.Msg)
}
