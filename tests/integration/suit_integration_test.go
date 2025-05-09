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
	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/timofey15g/homework/internal/handlers"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	pb "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

type AcceptOrderTestSuite struct {
	suite.Suite
	client  pb.OrderServiceClient
	handler *handlers.AcceptOrder
}

func (suite *AcceptOrderTestSuite) SetupTest() {
	suite.client = setupTest(suite.T())
	suite.handler = handlers.NewAcceptOrder(suite.client)
}

func (suite *AcceptOrderTestSuite) TestSuccessfulExecution() {
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

	suite.handler.Execute(w, req)

	resp := w.Result()
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *AcceptOrderTestSuite) TestInvalidRequestBody() {
	invalidBody := `{"id": "invalid_id"}`
	req := httptest.NewRequest(http.MethodPost, "/accept", bytes.NewReader([]byte(invalidBody)))
	w := httptest.NewRecorder()

	suite.handler.Execute(w, req)

	resp := w.Result()
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func TestAcceptOrderTestSuite(t *testing.T) {
	suite.Run(t, new(AcceptOrderTestSuite))
}

//////////////////////////////////////////////////

type IssueOrderTestSuite struct {
	suite.Suite
	client  pb.OrderServiceClient
	handler *handlers.IssueOrder
}

func (suite *IssueOrderTestSuite) SetupTest() {
	suite.client = setupTest(suite.T())
	suite.handler = handlers.NewIssueOrder(suite.client)
}

func (suite *IssueOrderTestSuite) TestSuccessfulExecution() {
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

		_, err := suite.client.CreateOrder(suite.T().Context(), req)
		assert.NoError(suite.T(), err)
	}

	requestBody, _ := json.Marshal([]int{1, 2, 3})
	req := httptest.NewRequest(http.MethodPost, "/orders/issue", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()

	suite.handler.Execute(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var orders models.OrdersSliceStorage
	err := json.NewDecoder(resp.Body).Decode(&orders)

	assert.NoError(suite.T(), err)
}

func (suite *IssueOrderTestSuite) TestInvalidRequestBody() {
	req := httptest.NewRequest(http.MethodPost, "/issue", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	suite.handler.Execute(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *IssueOrderTestSuite) TestStorageError() {
	requestBody, _ := json.Marshal([]int64{1, 2})
	req := httptest.NewRequest(http.MethodPost, "/issue", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()

	suite.handler.Execute(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.StatusCode)
}

func TestIssueOrderTestSuite(t *testing.T) {
	suite.Run(t, new(IssueOrderTestSuite))
}

////////////////////////////////////////////////////

type ListHistoryTestSuite struct {
	suite.Suite
	client  pb.OrderServiceClient
	handler *handlers.ListHistory
}

func (suite *ListHistoryTestSuite) SetupTest() {
	suite.client = setupTest(suite.T())
	suite.handler = handlers.NewListHistory(suite.client)

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

		_, err := suite.client.CreateOrder(suite.T().Context(), req)
		assert.NoError(suite.T(), err)
	}
}

func (suite *ListHistoryTestSuite) TestSuccessfulExecution() {
	req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
	w := httptest.NewRecorder()

	suite.handler.Execute(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var orders models.OrdersSliceStorage
	err := json.NewDecoder(w.Body).Decode(&orders)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(orders))
	assert.Equal(suite.T(), int64(1), orders[0].ID)
	assert.Equal(suite.T(), int64(2), orders[1].ID)
}

func (suite *ListHistoryTestSuite) TestMissingLimitParameter() {
	req := httptest.NewRequest(http.MethodGet, "/?offset=0", nil)
	w := httptest.NewRecorder()

	suite.handler.Execute(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *ListHistoryTestSuite) TestMissingOffsetParameter() {
	req := httptest.NewRequest(http.MethodGet, "/?limit=2", nil)
	w := httptest.NewRecorder()

	suite.handler.Execute(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *ListHistoryTestSuite) TestInvalidLimitParameter() {
	req := httptest.NewRequest(http.MethodGet, "/?limit=invalid&offset=0", nil)
	w := httptest.NewRecorder()

	suite.handler.Execute(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *ListHistoryTestSuite) TestInvalidOffsetParameter() {
	req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=invalid", nil)
	w := httptest.NewRecorder()

	suite.handler.Execute(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func TestListHistoryTestSuite(t *testing.T) {
	suite.Run(t, new(ListHistoryTestSuite))
}

/////////////////////////////////////////////////////////////

type ListOrderTestSuite struct {
	suite.Suite
	client  pb.OrderServiceClient
	handler *handlers.ListOrder
}

func (suite *ListOrderTestSuite) SetupTest() {
	suite.client = setupTest(suite.T())
	suite.handler = handlers.NewListOrder(suite.client)

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

		_, err := suite.client.CreateOrder(suite.T().Context(), req)

		assert.NoError(suite.T(), err)
	}
}

func (suite *ListOrderTestSuite) TestSuccessfulExecution() {
	req := httptest.NewRequest(http.MethodGet, "/orders?user_id=1&limit=10&cursor_id=0", nil)
	w := httptest.NewRecorder()

	suite.handler.Execute(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var orders models.OrdersSliceStorage
	err := json.NewDecoder(resp.Body).Decode(&orders)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), orders[0].ID)
	assert.Equal(suite.T(), int64(2), orders[1].ID)
}

func (suite *ListOrderTestSuite) TestInvalidQueryParameters() {
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
		suite.T().Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.query, nil)
			w := httptest.NewRecorder()

			suite.handler.Execute(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(suite.T(), tc.statusCode, resp.StatusCode)
		})
	}
}

func TestListOrderTestSuite(t *testing.T) {
	suite.Run(t, new(ListOrderTestSuite))
}

/////////////////////////////////////////////////////////

type ListReturnTestSuite struct {
	suite.Suite
	client  pb.OrderServiceClient
	handler *handlers.ListReturn
}

func (suite *ListReturnTestSuite) SetupTest() {
	suite.client = setupTest(suite.T())
	suite.handler = handlers.NewListReturn(suite.client)

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

		_, err := suite.client.CreateOrder(suite.T().Context(), req)
		assert.NoError(suite.T(), err)

		reqIssue := &pb.TReqIssueOrder{
			Ids: []int64{expectedOrders[i].ID},
		}
		_, err = suite.client.IssueOrder(suite.T().Context(), reqIssue)
		assert.NoError(suite.T(), err)

		reqReturn := &pb.TReqReturnOrder{
			OrderID: expectedOrders[i].ID,
			UserID:  expectedOrders[i].UserID,
		}

		_, err = suite.client.ReturnOrder(suite.T().Context(), reqReturn)
		assert.NoError(suite.T(), err)
	}
}

func (suite *ListReturnTestSuite) TestSuccessfulExecution() {
	req := httptest.NewRequest(http.MethodGet, "/returns?limit=10&offset=0", nil)
	w := httptest.NewRecorder()

	suite.handler.Execute(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var orders models.OrdersSliceStorage
	err := json.NewDecoder(resp.Body).Decode(&orders)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), orders[0].ID)
	assert.Equal(suite.T(), int64(2), orders[1].ID)
}

func (suite *ListReturnTestSuite) TestInvalidQueryParameters() {
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
		suite.T().Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.query, nil)
			w := httptest.NewRecorder()

			suite.handler.Execute(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(suite.T(), tc.statusCode, resp.StatusCode)
		})
	}
}

func TestListReturnTestSuite(t *testing.T) {
	suite.Run(t, new(ListReturnTestSuite))
}

//////////////////////////////////////////////////////////

type ReturnOrderTestSuite struct {
	suite.Suite
	client  pb.OrderServiceClient
	handler *handlers.ReturnOrder
}

func (suite *ReturnOrderTestSuite) SetupTest() {
	suite.client = setupTest(suite.T())
	suite.handler = handlers.NewReturnOrder(suite.client)
}

func (suite *ReturnOrderTestSuite) TestSuccessfulExecution() {
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

	_, err := suite.client.CreateOrder(suite.T().Context(), reqCreate)
	assert.NoError(suite.T(), err)

	reqIssue := &pb.TReqIssueOrder{
		Ids: []int64{expectedOrder.ID},
	}
	_, err = suite.client.IssueOrder(suite.T().Context(), reqIssue)
	assert.NoError(suite.T(), err)

	req := httptest.NewRequest(http.MethodGet, "/return?order_id=1&user_id=1", nil)
	w := httptest.NewRecorder()

	suite.handler.Execute(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var order models.Order
	err = json.NewDecoder(resp.Body).Decode(&order)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), order.ID, int64(1))
}

func (suite *ReturnOrderTestSuite) TestInvalidQueryParameters() {
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
		suite.T().Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.query, nil)
			w := httptest.NewRecorder()

			suite.handler.Execute(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(suite.T(), tc.statusCode, resp.StatusCode)
		})
	}
}

func (suite *ReturnOrderTestSuite) TestStorageError() {
	req := httptest.NewRequest(http.MethodGet, "/return?order_id=1&user_id=1", nil)
	w := httptest.NewRecorder()

	suite.handler.Execute(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.StatusCode)
}

func TestReturnOrderTestSuite(t *testing.T) {
	suite.Run(t, new(ReturnOrderTestSuite))
}

///////////////////////////////////////////////////////////

type WithdrawOrderTestSuite struct {
	suite.Suite
	client  pb.OrderServiceClient
	handler *handlers.WithdrawOrder
}

func (suite *WithdrawOrderTestSuite) SetupTest() {
	suite.client = setupTest(suite.T())
	suite.handler = handlers.NewWithdrawOrder(suite.client)
}

func (suite *WithdrawOrderTestSuite) TestSuccessfulExecution() {
	date := models.DefaultTime.Add(-480 * time.Hour)
	expectedOrder := models.NewOrder(1, 1, 10, date, 12.3, models.NewMoneyFromInt(100, 0), models.PackagingFilm, models.PackagingDefault)

	reqCreate := &pb.TReqAcceptOrder{
		ID:                  expectedOrder.ID,
		UserID:              expectedOrder.UserID,
		StorageDurationDays: 36500,
		Weight:              expectedOrder.Weight,
		Cost:                "100",
		Package:             string(expectedOrder.Package),
		ExtraPackage:        string(expectedOrder.ExtraPackage),
	}

	_, err := suite.client.CreateOrder(suite.T().Context(), reqCreate)
	assert.NoError(suite.T(), err)

	req := httptest.NewRequest(http.MethodGet, "/withdraw?order_id=1", nil)
	w := httptest.NewRecorder()

	suite.handler.Execute(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var order models.Order
	err = json.NewDecoder(resp.Body).Decode(&order)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedOrder.ID, order.ID)
}

func (suite *WithdrawOrderTestSuite) TestInvalidQueryParameters() {
	testCases := []struct {
		name       string
		query      string
		statusCode int
	}{
		{"Missing order_id", "/withdraw", http.StatusBadRequest},
		{"Invalid order_id", "/withdraw?order_id=abc", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.query, nil)
			w := httptest.NewRecorder()

			suite.handler.Execute(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(suite.T(), tc.statusCode, resp.StatusCode)
		})
	}
}

func (suite *WithdrawOrderTestSuite) TestStorageError() {
	req := httptest.NewRequest(http.MethodGet, "/withdraw?order_id=1", nil)
	w := httptest.NewRecorder()

	suite.handler.Execute(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.StatusCode)
}

func TestWithdrawOrderTestSuite(t *testing.T) {
	suite.Run(t, new(WithdrawOrderTestSuite))
}
