package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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

	t.Run("success", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := NewAcceptOrder(mockStorage, mock.NewMockLogPipeline())

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
		mockStorage := mock.NewMockStorage(ctrl)
		handler := NewAcceptOrder(mockStorage, mock.NewMockLogPipeline())

		req := httptest.NewRequest(http.MethodPost, "/accept", bytes.NewReader([]byte("invalid body")))
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "Invalid request body\n", rec.Body.String())
	})

	t.Run("storage error", func(t *testing.T) {
		mockStorage := mock.NewMockStorage(ctrl)
		handler := NewAcceptOrder(mockStorage, mock.NewMockLogPipeline())

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
