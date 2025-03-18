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

func TestListOrder_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := handlers.NewListOrder(mockStorage)

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
	})

	t.Run("bad request - missing user_id", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := handlers.NewListOrder(mockStorage)

		req := httptest.NewRequest(http.MethodGet, "/?limit=10&cursor_id=0", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := handlers.NewListOrder(mockStorage)

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
