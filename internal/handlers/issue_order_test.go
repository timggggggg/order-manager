package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/timofey15g/homework/internal/handlers/mock"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

func TestIssueOrder_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewIssueOrder(mockOrderServiceClient)

		ids := []int64{1, 2, 3}

		mockOrderServiceClient.EXPECT().
			IssueOrder(gomock.Any(), &pb.TReqIssueOrder{Ids: ids}).
			Return(&pb.TStringResp{Msg: "orders issued!"}, nil)

		body, _ := json.Marshal(ids)
		req := httptest.NewRequest(http.MethodPost, "/orders/issue", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "orders issued!", rec.Body.String())
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewIssueOrder(mockOrderServiceClient)

		req := httptest.NewRequest(http.MethodPost, "/orders/issue", bytes.NewReader([]byte("invalid body")))
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "Invalid request body\n", rec.Body.String())
	})

	t.Run("storage error", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewIssueOrder(mockOrderServiceClient)

		ids := []int64{1, 2, 3}

		mockOrderServiceClient.EXPECT().
			IssueOrder(gomock.Any(), &pb.TReqIssueOrder{Ids: ids}).
			Return(nil, errors.New("storage error"))

		body, _ := json.Marshal(ids)
		req := httptest.NewRequest(http.MethodPost, "/orders/issue", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "storage error\n", rec.Body.String())
	})
}
