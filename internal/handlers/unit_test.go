package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/timofey15g/homework/internal/handlers/mock"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/internal/packaging"
)

func TestAcceptOrder_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock.NewMockStorage(ctrl)
	handler := NewAcceptOrder(mockStorage)

	t.Run("success", func(t *testing.T) {
		orderJSON := OrderJSON{
			ID:                  1,
			UserID:              123,
			StorageDurationDays: 10,
			Weight:              5.0,
			Cost:                "100.00",
			Package:             "box",
			ExtraPackage:        "film",
		}

		body, _ := json.Marshal(orderJSON)
		req := httptest.NewRequest(http.MethodPost, "/accept", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		packagingStrategy, _ := packaging.NewPackagingStrategy(orderJSON.Package, packaging.PackagingStrategies)
		extraPackagingStrategy, _ := packaging.NewPackagingStrategy(orderJSON.ExtraPackage, packaging.ExtraPackagingStrategies)

		money, _ := models.NewMoney(orderJSON.Cost)
		acceptTime := time.Now()
		order := models.NewOrder(orderJSON.ID, orderJSON.UserID, orderJSON.StorageDurationDays, acceptTime,
			orderJSON.Weight, money, packagingStrategy.Type(), extraPackagingStrategy.Type())

		mockStorage.EXPECT().
			CreateOrder(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, o *models.Order) error {
				assert.Equal(t, order.ID, o.ID)
				assert.Equal(t, order.UserID, o.UserID)
				return nil
			})

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/accept", bytes.NewReader([]byte("invalid body")))
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "Invalid request body\n", rec.Body.String())
	})

	t.Run("storage error", func(t *testing.T) {

		orderJSON := OrderJSON{
			ID:                  1,
			UserID:              123,
			StorageDurationDays: 10,
			Weight:              5.0,
			Cost:                "100.00",
			Package:             "box",
			ExtraPackage:        "film",
		}

		body, _ := json.Marshal(orderJSON)
		req := httptest.NewRequest(http.MethodPost, "/accept", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		mockStorage.EXPECT().
			CreateOrder(gomock.Any(), gomock.Any()).
			Return(errors.New("storage error"))

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "error accepting order: storage error\n", rec.Body.String())
	})
}

func TestIssueOrder_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock.NewMockStorage(ctrl)
	handler := NewIssueOrder(mockStorage)
	t.Run("success", func(t *testing.T) {

		ids := []int64{1, 2, 3}

		expectedOrders := models.OrdersSliceStorage{
			models.NewOrder(1, 1, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(2, 2, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(3, 3, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
		}

		mockStorage.EXPECT().
			IssueOrders(gomock.Any(), ids).
			Return(expectedOrders, nil)

		body, _ := json.Marshal(ids)
		req := httptest.NewRequest(http.MethodPost, "/orders/issue", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var actualOrders models.OrdersSliceStorage
		err := json.NewDecoder(rec.Body).Decode(&actualOrders)
		assert.NoError(t, err)
		assert.Equal(t, expectedOrders, actualOrders)
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/orders/issue", bytes.NewReader([]byte("invalid body")))
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "Invalid request body\n", rec.Body.String())
	})

	t.Run("storage error", func(t *testing.T) {

		ids := []int64{1, 2, 3}
		mockStorage.EXPECT().
			IssueOrders(gomock.Any(), ids).
			Return(nil, errors.New("storage error"))

		body, _ := json.Marshal(ids)
		req := httptest.NewRequest(http.MethodPost, "/orders/issue", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "storage error\n", rec.Body.String())
	})
}

func TestListHistory_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock.NewMockStorage(ctrl)
	handler := NewListHistory(mockStorage)

	t.Run("success", func(t *testing.T) {
		limit := int64(10)
		offset := int64(0)
		expectedOrders := models.OrdersSliceStorage{
			models.NewOrder(1, 1, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(2, 2, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(3, 3, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
		}

		mockStorage.EXPECT().GetAll(gomock.Any(), limit, offset).Return(expectedOrders, nil)

		req := httptest.NewRequest(http.MethodGet, "/orders/?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10), nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var actualOrders models.OrdersSliceStorage
		err := json.NewDecoder(w.Result().Body).Decode(&actualOrders)
		assert.NoError(t, err)
		assert.Equal(t, expectedOrders, actualOrders)
	})

	t.Run("missing limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/orders/?offset=0", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("invalid limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/orders/?limit=abc&offset=0", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("storage error", func(t *testing.T) {
		limit := int64(10)
		offset := int64(0)
		mockStorage.EXPECT().GetAll(gomock.Any(), limit, offset).Return(nil, errors.New("storage error"))

		req := httptest.NewRequest(http.MethodGet, "/orders/?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10), nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}

func TestListOrder_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock.NewMockStorage(ctrl)
	handler := NewListOrder(mockStorage)

	t.Run("success", func(t *testing.T) {
		userID := int64(1)
		limit := int64(10)
		cursorID := int64(0)
		expectedOrders := models.OrdersSliceStorage{
			models.NewOrder(1, 1, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(2, 2, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(3, 3, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
		}

		mockStorage.EXPECT().
			GetByUserIDCursorPagination(gomock.Any(), userID, limit, cursorID).
			Return(expectedOrders, nil)

		req := httptest.NewRequest(http.MethodGet, "/?user_id="+strconv.FormatInt(userID, 10)+"&limit="+strconv.FormatInt(limit, 10)+"&cursor_id="+strconv.FormatInt(cursorID, 10), nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actualOrders models.OrdersSliceStorage
		err := json.NewDecoder(w.Body).Decode(&actualOrders)
		assert.NoError(t, err)
		assert.Equal(t, expectedOrders, actualOrders)
	})

	t.Run("bad request - missing user_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?limit=10&cursor_id=0", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		userID := int64(1)
		limit := int64(10)
		cursorID := int64(0)

		mockStorage.EXPECT().
			GetByUserIDCursorPagination(gomock.Any(), userID, limit, cursorID).
			Return(nil, errors.New("storage error"))

		req := httptest.NewRequest(http.MethodGet, "/?user_id="+strconv.FormatInt(userID, 10)+"&limit="+strconv.FormatInt(limit, 10)+"&cursor_id="+strconv.FormatInt(cursorID, 10), nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestListReturn_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock.NewMockStorage(ctrl)
	handler := NewListReturn(mockStorage)

	t.Run("success", func(t *testing.T) {
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
		assert.Equal(t, expectedOrders, actualOrders)
	})

	t.Run("missing limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?offset=5", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing offset", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?limit=10", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?limit=invalid&offset=5", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid offset", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?limit=10&offset=invalid", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("storage error", func(t *testing.T) {
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

func TestReturnOrder_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock.NewMockStorage(ctrl)
	handler := NewReturnOrder(mockStorage)

	t.Run("success", func(t *testing.T) {
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
		req := httptest.NewRequest(http.MethodGet, "/?user_id=456", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("missing user_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?order_id=123", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid order_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?order_id=abc&user_id=456", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid user_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?order_id=123&user_id=abc", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("storage error", func(t *testing.T) {
		orderID := int64(123)
		userID := int64(456)

		mockStorage.EXPECT().ReturnOrder(gomock.Any(), orderID, userID).Return(nil, errors.New("storage error"))

		req := httptest.NewRequest(http.MethodGet, "/?order_id="+strconv.FormatInt(orderID, 10)+"&user_id="+strconv.FormatInt(userID, 10), nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestWithdrawOrder_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock.NewMockStorage(ctrl)
	handler := NewWithdrawOrder(mockStorage)

	t.Run("success", func(t *testing.T) {
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
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid order_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?order_id=invalid", nil)
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("storage error", func(t *testing.T) {
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
