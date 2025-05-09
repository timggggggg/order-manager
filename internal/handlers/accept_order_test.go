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
	"google.golang.org/grpc"

	"gitlab.ozon.dev/timofey15g/homework/internal/handlers/mock"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/internal/packaging"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

func TestAcceptOrder_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewAcceptOrder(mockOrderServiceClient)

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

		mockOrderServiceClient.EXPECT().
			CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, r *pb.TReqAcceptOrder, opts ...grpc.CallOption) (*pb.TStringResp, error) {
				assert.Equal(t, order.ID, r.ID)
				assert.Equal(t, order.UserID, r.UserID)
				return &pb.TStringResp{Msg: "Order accepted successfully"}, nil
			})

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewAcceptOrder(mockOrderServiceClient)

		req := httptest.NewRequest(http.MethodPost, "/accept", bytes.NewReader([]byte("invalid body")))
		rec := httptest.NewRecorder()

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "Invalid request body\n", rec.Body.String())
	})

	t.Run("storage error", func(t *testing.T) {
		mockOrderServiceClient := mock.NewMockOrderServiceClient(ctrl)
		handler := NewAcceptOrder(mockOrderServiceClient)

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

		mockOrderServiceClient.EXPECT().
			CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, errors.New("error accepting order: storage error"))

		handler.Execute(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "error accepting order: storage error\n", rec.Body.String())
	})
}
