package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/timofey15g/homework/internal/handlers/mock"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

func TestReturnOrder_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := NewReturnOrder(mockStorage)

		orderID := int64(123)
		userID := int64(456)
		expectedOrder := &models.Order{ID: orderID, UserID: userID}

		mockStorage.EXPECT().ReturnOrder(gomock.Any(), orderID, userID).Return(expectedOrder, nil)

		req := httptest.NewRequest(http.MethodGet, "/?order_id="+strconv.FormatInt(orderID, 10)+"&user_id="+strconv.FormatInt(userID, 10), nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var actualOrder models.Order
		err := json.NewDecoder(rec.Body).Decode(&actualOrder)
		assert.NoError(t, err)
		assert.Equal(t, *expectedOrder, actualOrder)
	})

	t.Run("missing order_id", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := NewReturnOrder(mockStorage)

		req := httptest.NewRequest(http.MethodGet, "/?user_id=456", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("missing user_id", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := NewReturnOrder(mockStorage)

		req := httptest.NewRequest(http.MethodGet, "/?order_id=123", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid order_id", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := NewReturnOrder(mockStorage)

		req := httptest.NewRequest(http.MethodGet, "/?order_id=abc&user_id=456", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid user_id", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := NewReturnOrder(mockStorage)

		req := httptest.NewRequest(http.MethodGet, "/?order_id=123&user_id=abc", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("storage error", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := NewReturnOrder(mockStorage)

		orderID := int64(123)
		userID := int64(456)

		mockStorage.EXPECT().ReturnOrder(gomock.Any(), orderID, userID).Return(nil, errors.New("storage error"))

		req := httptest.NewRequest(http.MethodGet, "/?order_id="+strconv.FormatInt(orderID, 10)+"&user_id="+strconv.FormatInt(userID, 10), nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
