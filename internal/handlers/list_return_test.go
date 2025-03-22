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

func TestListReturn_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := NewListReturn(mockStorage, mock.NewMockLogPipeline())

		limit := int64(10)
		offset := int64(5)
		expectedOrders := models.OrdersSliceStorage{
			models.NewOrder(1, 1, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(2, 2, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(3, 3, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
		}

		mockStorage.
			EXPECT().
			GetReturnsLimitOffsetPagination(gomock.Any(), limit, offset).
			Return(expectedOrders, nil)

		req := httptest.NewRequest(http.MethodGet, "/?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10), nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actualOrders models.OrdersSliceStorage
		err := json.NewDecoder(w.Body).Decode(&actualOrders)
		assert.NoError(t, err)
	})

	t.Run("missing limit", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := NewListReturn(mockStorage, mock.NewMockLogPipeline())

		req := httptest.NewRequest(http.MethodGet, "/?offset=5", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing offset", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := NewListReturn(mockStorage, mock.NewMockLogPipeline())

		req := httptest.NewRequest(http.MethodGet, "/?limit=10", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid limit", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := NewListReturn(mockStorage, mock.NewMockLogPipeline())

		req := httptest.NewRequest(http.MethodGet, "/?limit=invalid&offset=5", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid offset", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := NewListReturn(mockStorage, mock.NewMockLogPipeline())

		req := httptest.NewRequest(http.MethodGet, "/?limit=10&offset=invalid", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("storage error", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := NewListReturn(mockStorage, mock.NewMockLogPipeline())

		limit := int64(10)
		offset := int64(5)
		expectedError := errors.New("storage error")

		mockStorage.
			EXPECT().
			GetReturnsLimitOffsetPagination(gomock.Any(), limit, offset).
			Return(nil, expectedError)

		req := httptest.NewRequest(http.MethodGet, "/?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10), nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
