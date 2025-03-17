package unit

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/timofey15g/homework/internal/handlers"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/tests/unit/mock"
)

func TestWithdrawOrder_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := handlers.NewWithdrawOrder(mockStorage)

		orderID := int64(123)
		expectedOrder := &models.Order{
			ID: orderID,
		}

		mockStorage.EXPECT().
			WithdrawOrder(gomock.Any(), orderID).
			Return(expectedOrder, nil)

		req := httptest.NewRequest(http.MethodGet, "/?order_id="+strconv.FormatInt(orderID, 10), nil)
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
		handler := handlers.NewWithdrawOrder(mockStorage)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid order_id", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := handlers.NewWithdrawOrder(mockStorage)

		req := httptest.NewRequest(http.MethodGet, "/?order_id=invalid", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("storage error", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := handlers.NewWithdrawOrder(mockStorage)

		orderID := int64(123)
		mockStorage.EXPECT().
			WithdrawOrder(gomock.Any(), orderID).
			Return(nil, errors.New("storage error"))

		req := httptest.NewRequest(http.MethodGet, "/?order_id="+strconv.FormatInt(orderID, 10), nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
