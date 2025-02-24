package storage

import (
	"encoding/json"
	"os"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type Storage struct {
	orders   []*models.Order
	filePath string
}

func NewStorage(filePath string) (*Storage, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	orders := make([]*models.Order, 0)
	err = json.Unmarshal(file, &orders)
	if err != nil {
		return nil, err
	}

	st := &Storage{
		orders,
		filePath,
	}

	return st, nil
}

func (s *Storage) Save() error {
	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	orders, err := json.MarshalIndent(s.orders, "", "\t")
	if err != nil {
		return err
	}

	_, err = file.Write(orders)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetByID(id int64) (*models.Order, error) {
	for _, order := range s.orders {
		if order.ID == id {
			return order, nil
		}
	}

	return nil, models.ErrorOrderNotFound
}

func (s *Storage) DeleteByID(id int64) error {
	for i, order := range s.orders {
		if order.ID == id {
			s.orders = append(s.orders[:i], s.orders[i+1:]...)
			return s.Save()
		}
	}

	return models.ErrorOrderNotFound
}

func (s *Storage) Update(updatedOrder *models.Order) error {
	for i, order := range s.orders {
		if order.ID == updatedOrder.ID {
			s.orders[i] = updatedOrder
			return s.Save()
		}
	}

	return models.ErrorOrderNotFound
}

func (s *Storage) Add(order *models.Order) error {
	exists, err := s.GetByID(order.ID)
	if err == nil && exists != nil {
		return models.ErrorOrderAlreadyExists
	}

	s.orders = append(s.orders, order)

	return s.Save()
}

func (s *Storage) GetAllOrders() []*models.Order {
	return s.orders
}

func (s *Storage) GetSize() int64 {
	return int64(len(s.orders))
}
