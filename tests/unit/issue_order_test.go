package unit

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/timofey15g/homework/internal/handlers"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/tests/unit/mock"
)

func TestIssueOrder_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := handlers.NewIssueOrder(mockStorage)

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
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := handlers.NewIssueOrder(mockStorage)

		req := httptest.NewRequest(http.MethodPost, "/orders/issue", bytes.NewReader([]byte("invalid body")))
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "Invalid request body\n", rec.Body.String())
	})

	t.Run("storage error", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := handlers.NewIssueOrder(mockStorage)

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
