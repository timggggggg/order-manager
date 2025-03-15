// nolint
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/internal/storage/postgres"
)

func newPgFacade(pool *pgxpool.Pool) *postgres.PgFacade {
	txManager := postgres.NewTxManager(pool)
	pgRepository := postgres.NewPgRepository(txManager)
	return postgres.NewPgFacade(txManager, pgRepository)
}

func newPgxPool(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func setupTest(t *testing.T) *postgres.PgFacade {
	err := godotenv.Load(".env")
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

	return newPgFacade(pool)
}

func TestAcceptOrder_Execute_integration(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		storage := setupTest(t)

		handler := NewAcceptOrder(storage)

		orderJSON := OrderJSON{
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
		handler := NewAcceptOrder(storage)

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
		handler := NewIssueOrder(storage)
		expectedOrders := models.OrdersSliceStorage{
			models.NewOrder(1, 1, 10, time.Now(), 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(2, 1, 10, time.Now(), 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
			models.NewOrder(3, 1, 10, time.Now(), 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault),
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

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading the response body:", err)
			return
		}

		// Print the raw response body
		fmt.Println("Raw Response Body:")
		fmt.Println(body)

		assert.NoError(t, err)
	})

	t.Run("invalid request body", func(t *testing.T) {
		storage := setupTest(t)
		handler := NewIssueOrder(storage)

		req := httptest.NewRequest(http.MethodPost, "/issue", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		handler.Execute(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("storage error", func(t *testing.T) {
		storage := setupTest(t)
		handler := NewIssueOrder(storage)

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

}

func TestListOrder_Execute_integration(t *testing.T) {

}

func TestListReturn_Execute_integration(t *testing.T) {

}

func TestReturnOrder_Execute_integration(t *testing.T) {

}

func TestWithdrawOrder_Execute_integration(t *testing.T) {

}
