// nolint
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/timofey15g/homework/internal/handlers"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/internal/storage/postgres"
	storagecache "gitlab.ozon.dev/timofey15g/homework/internal/storage_cache"
)

func newPgFacade(t *testing.T, pool *pgxpool.Pool) *postgres.PgFacade {
	txManager := postgres.NewTxManager(pool)
	pgRepository := postgres.NewPgRepository(txManager)
	cache, err := storagecache.NewCacheStrategy("DEFAULT", int64(0))
	if err != nil {
		t.Fatal("Error while creating strategy")
	}
	return postgres.NewPgFacade(txManager, pgRepository, cache, time.Now)
}

func newPgxPool(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func setupTest(t *testing.T) *postgres.PgFacade {
	err := godotenv.Load(".env.example")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	db := os.Getenv("POSTGRES_DB")
	sslMode := os.Getenv("SSL_MODE")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, db, sslMode)

	pool, err := newPgxPool(ctx, connectionString)
	if err != nil {
		t.Fatal("error newPgxPool", err)
	}

	t.Cleanup(func() {
		pool.Exec(ctx, "TRUNCATE TABLE orders RESTART IDENTITY CASCADE")
		pool.Close()
	})

	return newPgFacade(t, pool)
}

func TestAcceptOrder_Execute_integration(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		storage := setupTest(t)

		handler := handlers.NewAcceptOrder(storage)

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
		storage := setupTest(t)
		handler := handlers.NewAcceptOrder(storage)

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
		storage := setupTest(t)
		handler := handlers.NewIssueOrder(storage)
		expectedOrders := models.OrdersSliceStorage{
			models.NewOrder(1, 1, 36500, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(2, 1, 36500, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(3, 1, 36500, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
		}

		for i := range expectedOrders {
			err := storage.CreateOrder(t.Context(), expectedOrders[i])
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
		storage := setupTest(t)
		handler := handlers.NewIssueOrder(storage)

		req := httptest.NewRequest(http.MethodPost, "/issue", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("storage error", func(t *testing.T) {
		storage := setupTest(t)
		handler := handlers.NewIssueOrder(storage)

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
	storage := setupTest(t)
	handler := handlers.NewListHistory(storage)

	expectedOrders := models.OrdersSliceStorage{
		models.NewOrder(2, 1, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
		models.NewOrder(1, 1, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
	}

	for i := range expectedOrders {
		err := storage.CreateOrder(t.Context(), expectedOrders[i])
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
		storage := setupTest(t)
		handler := handlers.NewListOrder(storage)

		expectedOrders := models.OrdersSliceStorage{
			models.NewOrder(1, 1, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(2, 1, 10, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
		}

		for i := range expectedOrders {
			err := storage.CreateOrder(t.Context(), expectedOrders[i])
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
		storage := setupTest(t)
		handler := handlers.NewListOrder(storage)

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
		storage := setupTest(t)
		handler := handlers.NewListReturn(storage)

		expectedOrders := models.OrdersSliceStorage{
			models.NewOrder(1, 1, 36500, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(2, 1, 36500, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
		}

		for i := range expectedOrders {
			err := storage.CreateOrder(t.Context(), expectedOrders[i])
			assert.NoError(t, err)

			_, err = storage.IssueOrders(t.Context(), []int64{expectedOrders[i].ID})
			assert.NoError(t, err)

			_, err = storage.ReturnOrder(t.Context(), expectedOrders[i].ID, expectedOrders[i].UserID)
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
		storage := setupTest(t)
		handler := handlers.NewListReturn(storage)

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
		storage := setupTest(t)
		handler := handlers.NewReturnOrder(storage)

		expectedOrder := models.NewOrder(1, 1, 36500, models.DefaultTime, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault)

		err := storage.CreateOrder(t.Context(), expectedOrder)
		assert.NoError(t, err)
		_, err = storage.IssueOrders(t.Context(), []int64{1})
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
		storage := setupTest(t)
		handler := handlers.NewReturnOrder(storage)

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

	t.Run("storage error", func(t *testing.T) {
		storage := setupTest(t)
		handler := handlers.NewReturnOrder(storage)

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
		storage := setupTest(t)
		handler := handlers.NewWithdrawOrder(storage)

		date := models.DefaultTime.Add(-480 * time.Hour)
		expectedOrder := models.NewOrder(1, 1, 10, date, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault)

		err := storage.CreateOrder(t.Context(), expectedOrder)
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
		storage := setupTest(t)
		handler := handlers.NewWithdrawOrder(storage)

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

	t.Run("storage error", func(t *testing.T) {
		storage := setupTest(t)
		handler := handlers.NewWithdrawOrder(storage)

		req := httptest.NewRequest(http.MethodGet, "/withdraw?order_id=1", nil)
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
