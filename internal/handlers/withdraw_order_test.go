package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/timofey15g/homework/internal/handlers/mock"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

func TestWithdrawOrder_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewWithdrawOrder(mockOrderServiceClient)

		orderID := int64(123)

		mockOrderServiceClient.EXPECT().
			WithdrawOrder(gomock.Any(), &pb.TReqWithdrawOrder{OrderID: orderID}).
			Return(&pb.TStringResp{Msg: ""}, nil)

		req := httptest.NewRequest(http.MethodGet, "/?order_id="+strconv.FormatInt(orderID, 10), nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Equal(t, rec.Body.String(), "")
	})

	t.Run("missing order_id", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewWithdrawOrder(mockOrderServiceClient)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid order_id", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewWithdrawOrder(mockOrderServiceClient)

		req := httptest.NewRequest(http.MethodGet, "/?order_id=invalid", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("storage error", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewWithdrawOrder(mockOrderServiceClient)

		orderID := int64(123)
		mockOrderServiceClient.EXPECT().
			WithdrawOrder(gomock.Any(), &pb.TReqWithdrawOrder{OrderID: orderID}).
			Return(nil, errors.New("storage error"))

		req := httptest.NewRequest(http.MethodGet, "/?order_id="+strconv.FormatInt(orderID, 10), nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
