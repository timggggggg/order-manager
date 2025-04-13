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

func TestListHistory_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewListHistory(mockOrderServiceClient)

		limit := int64(10)
		offset := int64(0)
		expectedOrders := models.OrdersSliceStorage{
			models.NewOrder(1, 1, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(2, 2, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(3, 3, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
		}

		mockOrderServiceClient.EXPECT().
			ListHistory(gomock.Any(), &pb.TReqListHistory{Limit: limit, Offset: offset}).
			Return(&pb.TListResp{Orders: models.OrdersSliceModelToProto(expectedOrders)}, nil)

		req := httptest.NewRequest(http.MethodGet, "/orders/?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10), nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var actualOrders models.OrdersSliceStorage
		err := json.NewDecoder(w.Result().Body).Decode(&actualOrders)
		assert.NoError(t, err)
	})

	t.Run("missing limit", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewListHistory(mockOrderServiceClient)

		req := httptest.NewRequest(http.MethodGet, "/orders/?offset=0", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("invalid limit", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewListHistory(mockOrderServiceClient)

		req := httptest.NewRequest(http.MethodGet, "/orders/?limit=abc&offset=0", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("storage error", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewListHistory(mockOrderServiceClient)

		limit := int64(10)
		offset := int64(0)
		mockOrderServiceClient.EXPECT().
			ListHistory(gomock.Any(), &pb.TReqListHistory{Limit: limit, Offset: offset}).
			Return(nil, errors.New("storage error"))

		req := httptest.NewRequest(http.MethodGet, "/orders/?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10), nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}
