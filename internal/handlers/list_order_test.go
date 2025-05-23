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
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

func TestListOrder_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewListOrder(mockOrderServiceClient)

		userID := int64(1)
		limit := int64(10)
		cursorID := int64(0)
		expectedOrders := models.OrdersSliceStorage{
			models.NewOrder(1, 1, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(2, 2, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(3, 3, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
		}

		mockOrderServiceClient.EXPECT().
			ListOrders(gomock.Any(), &pb.TReqListOrders{UserID: userID, Limit: limit, CursorID: cursorID}).
			Return(&pb.TListResp{Orders: models.OrdersSliceModelToProto(expectedOrders)}, nil)

		req := httptest.NewRequest(http.MethodGet, "/?user_id="+strconv.FormatInt(userID, 10)+"&limit="+strconv.FormatInt(limit, 10)+"&cursor_id="+strconv.FormatInt(cursorID, 10), nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actualOrders models.OrdersSliceStorage
		err := json.NewDecoder(w.Body).Decode(&actualOrders)
		assert.NoError(t, err)
	})

	t.Run("bad request - missing user_id", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewListOrder(mockOrderServiceClient)

		req := httptest.NewRequest(http.MethodGet, "/?limit=10&cursor_id=0", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewListOrder(mockOrderServiceClient)

		userID := int64(1)
		limit := int64(10)
		cursorID := int64(0)

		mockOrderServiceClient.EXPECT().
			ListOrders(gomock.Any(), &pb.TReqListOrders{UserID: userID, Limit: limit, CursorID: cursorID}).
			Return(nil, errors.New("storage error"))

		req := httptest.NewRequest(http.MethodGet, "/?user_id="+strconv.FormatInt(userID, 10)+"&limit="+strconv.FormatInt(limit, 10)+"&cursor_id="+strconv.FormatInt(cursorID, 10), nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
