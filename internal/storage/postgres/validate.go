package postgres

import (
	"time"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

func validateReturn(order *OrderDB, userID int64) error {
	if order.UserID != userID {
		return models.ErrorOrderInvalidUser
	}
	if err := validateReturnOrderStatus(order); err != nil {
		return err
	}
	if err := validateReturnExpired(order); err != nil {
		return err
	}
	order.Status = "returned"
	return nil
}

func validateReturnExpired(order *OrderDB) error {
	if returnDeadline := order.IssueTime.Time.Add(models.MaxReturnTime); returnDeadline.Before(time.Now()) {
		order.Status = "expired"
		return models.ErrorOrderReturnExpired
	}
	return nil
}

func validateReturnOrderStatus(order *OrderDB) error {
	if order.Status == string(models.StatusReturned) {
		return models.ErrorOrderAlreadyReturned
	}
	if order.Status == string(models.StatusAccepted) {
		return models.ErrorOrderNotIssued
	}
	return nil
}

func validateIssues(ordersMap OrdersDBMapStorage) error {
	if err := validateSameUser(ordersMap); err != nil {
		return err
	}
	for _, order := range ordersMap {
		if err := validateIssue(order); err != nil {
			return err
		}
		order.Status = "issued"
	}
	return nil
}

func validateSameUser(ordersMap OrdersDBMapStorage) error {
	userIDMap := make(map[int64]bool)
	for _, order := range ordersMap {
		userIDMap[order.UserID] = true
	}
	if len(userIDMap) > 1 {
		return models.ErrorOrdersDifferentUsers
	}
	return nil
}

func validateIssue(order *OrderDB) error {
	if order.Status == string(models.StatusReturned) {
		return models.ErrorOrderAlreadyReturned
	}
	if order.Status == string(models.StatusIssued) {
		return models.ErrorOrderAlreadyIssued
	}
	if order.ExpireTime.Time.Before(time.Now()) {
		return models.ErrorOrderStorageExpired
	}
	return nil
}

func validateWithdraw(order *OrderDB) error {
	if order.Status == string(models.StatusAccepted) && time.Now().Before(order.ExpireTime.Time) {
		return models.ErrorOrderNotExpired
	}
	order.Status = "withdrawed"
	return nil
}
