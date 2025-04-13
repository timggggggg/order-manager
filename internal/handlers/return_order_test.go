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

func TestReturnOrder_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewReturnOrder(mockOrderServiceClient)

		orderID := int64(123)
		userID := int64(456)

		mockOrderServiceClient.EXPECT().
			ReturnOrder(gomock.Any(), &pb.TReqReturnOrder{OrderID: orderID, UserID: userID}).
			Return(&pb.TStringResp{Msg: ""}, nil)

		req := httptest.NewRequest(http.MethodGet, "/?order_id="+strconv.FormatInt(orderID, 10)+"&user_id="+strconv.FormatInt(userID, 10), nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Equal(t, rec.Body.String(), "")
	})

	t.Run("missing order_id", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewReturnOrder(mockOrderServiceClient)

		req := httptest.NewRequest(http.MethodGet, "/?user_id=456", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("missing user_id", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewReturnOrder(mockOrderServiceClient)

		req := httptest.NewRequest(http.MethodGet, "/?order_id=123", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid order_id", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewReturnOrder(mockOrderServiceClient)

		req := httptest.NewRequest(http.MethodGet, "/?order_id=abc&user_id=456", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid user_id", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewReturnOrder(mockOrderServiceClient)

		req := httptest.NewRequest(http.MethodGet, "/?order_id=123&user_id=abc", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("storage error", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewReturnOrder(mockOrderServiceClient)

		orderID := int64(123)
		userID := int64(456)

		mockOrderServiceClient.EXPECT().
			ReturnOrder(gomock.Any(), &pb.TReqReturnOrder{OrderID: orderID, UserID: userID}).
			Return(nil, errors.New("storage error"))

		req := httptest.NewRequest(http.MethodGet, "/?order_id="+strconv.FormatInt(orderID, 10)+"&user_id="+strconv.FormatInt(userID, 10), nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
