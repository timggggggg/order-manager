// nolint
package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gitlab.ozon.dev/timofey15g/homework/internal/handlers"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

func setupTest(t *testing.T) pb.OrderServiceClient {
	conn, err := grpc.NewClient("localhost:5252", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to create grpc client: %v", err)
	}

	clnt := pb.NewOrderServiceClient(conn)

	return clnt
}

func TestAcceptOrder_Execute_integration(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		client := setupTest(t)

		handler := handlers.NewAcceptOrder(client)

		orderJSON := handlers.OrderJSON{
			ID:                  1,
			UserID:              123,
			StorageDurationDays: 10,
			Weight:              5.5,
			Cost:                "100.00",
			Package:             "box",
			ExtraPackage:        "film",
		}
		body, _ := json.Marshal(orderJSON)

		req := httptest.NewRequest(http.MethodPost, "/accept", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("invalid request body", func(t *testing.T) {
		client := setupTest(t)
		handler := handlers.NewAcceptOrder(client)

		invalidBody := `{"id": "invalid_id"}`
		req := httptest.NewRequest(http.MethodPost, "/accept", bytes.NewReader([]byte(invalidBody)))
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestIssueOrder_Execute_integration(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		client := setupTest(t)
		handler := handlers.NewIssueOrder(client)
		expectedOrders := models.OrdersSliceStorage{
			models.NewOrder(1, 1, 36500, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(2, 1, 36500, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(3, 1, 36500, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
		}

		for i := range expectedOrders {
			req := &pb.TReqAcceptOrder{
				ID:                  expectedOrders[i].ID,
				UserID:              expectedOrders[i].UserID,
				StorageDurationDays: 36500,
				Weight:              expectedOrders[i].Weight,
				Cost:                "100",
				Package:             string(expectedOrders[i].Package),
				ExtraPackage:        string(expectedOrders[i].ExtraPackage),
			}

			_, err := client.CreateOrder(t.Context(), req)
			assert.NoError(t, err)
		}

		requestBody, _ := json.Marshal([]int{1, 2, 3})
		req := httptest.NewRequest(http.MethodPost, "/orders/issue", bytes.NewReader(requestBody))
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var orders models.OrdersSliceStorage
		err := json.NewDecoder(resp.Body).Decode(&orders)

		assert.NoError(t, err)
	})

	t.Run("invalid request body", func(t *testing.T) {
		client := setupTest(t)
		handler := handlers.NewIssueOrder(client)

		req := httptest.NewRequest(http.MethodPost, "/issue", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("client error", func(t *testing.T) {
		client := setupTest(t)
		handler := handlers.NewIssueOrder(client)

		requestBody, _ := json.Marshal([]int64{1, 2})
		req := httptest.NewRequest(http.MethodPost, "/issue", bytes.NewReader(requestBody))
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestListHistory_Execute_integration(t *testing.T) {
	client := setupTest(t)
	handler := handlers.NewListHistory(client)

	expectedOrders := models.OrdersSliceStorage{
		models.NewOrder(2, 1, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
		models.NewOrder(1, 1, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
	}

	for i := range expectedOrders {
		req := &pb.TReqAcceptOrder{
			ID:                  expectedOrders[i].ID,
			UserID:              expectedOrders[i].UserID,
			StorageDurationDays: 10,
			Weight:              expectedOrders[i].Weight,
			Cost:                "100",
			Package:             string(expectedOrders[i].Package),
			ExtraPackage:        string(expectedOrders[i].ExtraPackage),
		}

		_, err := client.CreateOrder(t.Context(), req)
		assert.NoError(t, err)
	}

	t.Run("successful execution", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var orders models.OrdersSliceStorage
		err := json.NewDecoder(w.Body).Decode(&orders)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(orders))
		assert.Equal(t, int64(1), orders[0].ID)
		assert.Equal(t, int64(2), orders[1].ID)
	})

	t.Run("error due to missing limit parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?offset=0", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("error due to missing offset parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?limit=2", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("error due to invalid limit parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?limit=invalid&offset=0", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("error due to invalid offset parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=invalid", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestListOrder_Execute_integration(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		client := setupTest(t)
		handler := handlers.NewListOrder(client)

		expectedOrders := models.OrdersSliceStorage{
			models.NewOrder(1, 1, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(2, 1, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
		}

		for i := range expectedOrders {
			req := &pb.TReqAcceptOrder{
				ID:                  expectedOrders[i].ID,
				UserID:              expectedOrders[i].UserID,
				StorageDurationDays: 10,
				Weight:              expectedOrders[i].Weight,
				Cost:                "100",
				Package:             string(expectedOrders[i].Package),
				ExtraPackage:        string(expectedOrders[i].ExtraPackage),
			}

			_, err := client.CreateOrder(t.Context(), req)
			assert.NoError(t, err)
		}

		req := httptest.NewRequest(http.MethodGet, "/orders?user_id=1&limit=10&cursor_id=0", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var orders models.OrdersSliceStorage
		err := json.NewDecoder(resp.Body).Decode(&orders)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), orders[0].ID)
		assert.Equal(t, int64(2), orders[1].ID)
	})

	t.Run("invalid query parameters", func(t *testing.T) {
		client := setupTest(t)
		handler := handlers.NewListOrder(client)

		testCases := []struct {
			name       string
			query      string
			statusCode int
		}{
			{"Missing user_id", "/orders?limit=10&cursor_id=0", http.StatusBadRequest},
			{"Missing limit", "/orders?user_id=1&cursor_id=0", http.StatusBadRequest},
			{"Missing cursor_id", "/orders?user_id=1&limit=10", http.StatusBadRequest},
			{"Invalid user_id", "/orders?user_id=abc&limit=10&cursor_id=0", http.StatusBadRequest},
			{"Invalid limit", "/orders?user_id=1&limit=abc&cursor_id=0", http.StatusBadRequest},
			{"Invalid cursor_id", "/orders?user_id=1&limit=10&cursor_id=abc", http.StatusBadRequest},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, tc.query, nil)
				w := httptest.NewRecorder()

				handler.Execute(w, req)

				resp := w.Result()
				defer resp.Body.Close()

				assert.Equal(t, tc.statusCode, resp.StatusCode)
			})
		}
	})
}

func TestListReturn_Execute_integration(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		client := setupTest(t)
		handler := handlers.NewListReturn(client)

		expectedOrders := models.OrdersSliceStorage{
			models.NewOrder(1, 1, 36500, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(2, 1, 36500, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
		}

		for i := range expectedOrders {
			req := &pb.TReqAcceptOrder{
				ID:                  expectedOrders[i].ID,
				UserID:              expectedOrders[i].UserID,
				StorageDurationDays: 36500,
				Weight:              expectedOrders[i].Weight,
				Cost:                "100",
				Package:             string(expectedOrders[i].Package),
				ExtraPackage:        string(expectedOrders[i].ExtraPackage),
			}

			_, err := client.CreateOrder(t.Context(), req)
			assert.NoError(t, err)

			reqIssue := &pb.TReqIssueOrder{
				Ids: []int64{expectedOrders[i].ID},
			}
			_, err = client.IssueOrder(t.Context(), reqIssue)
			assert.NoError(t, err)

			reqReturn := &pb.TReqReturnOrder{
				OrderID: expectedOrders[i].ID,
				UserID:  expectedOrders[i].UserID,
			}

			_, err = client.ReturnOrder(t.Context(), reqReturn)
			assert.NoError(t, err)
		}

		req := httptest.NewRequest(http.MethodGet, "/returns?limit=10&offset=0", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var orders models.OrdersSliceStorage
		err := json.NewDecoder(resp.Body).Decode(&orders)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), orders[0].ID)
		assert.Equal(t, int64(2), orders[1].ID)
	})

	t.Run("invalid query parameters", func(t *testing.T) {
		client := setupTest(t)
		handler := handlers.NewListReturn(client)

		testCases := []struct {
			name       string
			query      string
			statusCode int
		}{
			{"Missing limit", "/returns?offset=0", http.StatusBadRequest},
			{"Missing offset", "/returns?limit=10", http.StatusBadRequest},
			{"Invalid limit", "/returns?limit=abc&offset=0", http.StatusBadRequest},
			{"Invalid offset", "/returns?limit=10&offset=abc", http.StatusBadRequest},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, tc.query, nil)
				w := httptest.NewRecorder()

				handler.Execute(w, req)

				resp := w.Result()
				defer resp.Body.Close()

				assert.Equal(t, tc.statusCode, resp.StatusCode)
			})
		}
	})
}

func TestReturnOrder_Execute_integration(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		client := setupTest(t)
		handler := handlers.NewReturnOrder(client)

		expectedOrder := models.NewOrder(1, 1, 36500, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault)

		reqCreate := &pb.TReqAcceptOrder{
			ID:                  expectedOrder.ID,
			UserID:              expectedOrder.UserID,
			StorageDurationDays: 36500,
			Weight:              expectedOrder.Weight,
			Cost:                "100",
			Package:             string(expectedOrder.Package),
			ExtraPackage:        string(expectedOrder.ExtraPackage),
		}

		_, err := client.CreateOrder(t.Context(), reqCreate)
		assert.NoError(t, err)

		reqIssue := &pb.TReqIssueOrder{
			Ids: []int64{expectedOrder.ID},
		}
		_, err = client.IssueOrder(t.Context(), reqIssue)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/return?order_id=1&user_id=1", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var order models.Order
		err = json.NewDecoder(resp.Body).Decode(&order)

		assert.NoError(t, err)
		assert.Equal(t, order.ID, int64(1))
	})

	t.Run("invalid query parameters", func(t *testing.T) {
		client := setupTest(t)
		handler := handlers.NewReturnOrder(client)

		testCases := []struct {
			name       string
			query      string
			statusCode int
		}{
			{"Missing order_id", "/return?user_id=1", http.StatusBadRequest},
			{"Missing user_id", "/return?order_id=1", http.StatusBadRequest},
			{"Invalid order_id", "/return?order_id=abc&user_id=1", http.StatusBadRequest},
			{"Invalid user_id", "/return?order_id=1&user_id=abc", http.StatusBadRequest},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, tc.query, nil)
				w := httptest.NewRecorder()

				handler.Execute(w, req)

				resp := w.Result()
				defer resp.Body.Close()

				assert.Equal(t, tc.statusCode, resp.StatusCode)
			})
		}
	})

	t.Run("client error", func(t *testing.T) {
		client := setupTest(t)
		handler := handlers.NewReturnOrder(client)

		req := httptest.NewRequest(http.MethodGet, "/return?order_id=1&user_id=1", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestWithdrawOrder_Execute_integration(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		client := setupTest(t)
		handler := handlers.NewWithdrawOrder(client)

		date := models.DefaultTime.Add(-480 * time.Hour)
		expectedOrder := models.NewOrder(1, 1, 10, date, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault)

		reqCreate := &pb.TReqAcceptOrder{
			ID:                  expectedOrder.ID,
			UserID:              expectedOrder.UserID,
			StorageDurationDays: 10,
			Weight:              expectedOrder.Weight,
			Cost:                "100",
			Package:             string(expectedOrder.Package),
			ExtraPackage:        string(expectedOrder.ExtraPackage),
		}

		_, err := client.CreateOrder(t.Context(), reqCreate)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/withdraw?order_id=1", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var order models.Order
		err = json.NewDecoder(resp.Body).Decode(&order)

		assert.NoError(t, err)
		assert.Equal(t, expectedOrder.ID, order.ID)
	})

	t.Run("invalid query parameters", func(t *testing.T) {
		client := setupTest(t)
		handler := handlers.NewWithdrawOrder(client)

		testCases := []struct {
			name       string
			query      string
			statusCode int
		}{
			{"Missing order_id", "/withdraw", http.StatusBadRequest},
			{"Invalid order_id", "/withdraw?order_id=abc", http.StatusBadRequest},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, tc.query, nil)
				w := httptest.NewRecorder()

				handler.Execute(w, req)

				resp := w.Result()
				defer resp.Body.Close()

				assert.Equal(t, tc.statusCode, resp.StatusCode)
			})
		}
	})

	t.Run("client error", func(t *testing.T) {
		client := setupTest(t)
		handler := handlers.NewWithdrawOrder(client)

		req := httptest.NewRequest(http.MethodGet, "/withdraw?order_id=1", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
