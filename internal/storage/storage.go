package storage

import (
	"encoding/json"
	"os"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type Storage struct {
	Orders   []*Order
	filePath string
}

func NewStorage(filePath string) (*Storage, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	orders := make([]*Order, 0)
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

	orders, err := json.MarshalIndent(s.Orders, "", "\t")
	if err != nil {
		return err
	}

	_, err = file.Write(orders)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetByID(id int64) (*Order, error) {
	for _, order := range s.Orders {
		if order.ID == id {
			return order, nil
		}
	}

	return nil, models.ErrorOrderNotFound
}

func (s *Storage) DeleteByID(id int64) error {
	for i, order := range s.Orders {
		if order.ID == id {
			s.Orders = append(s.Orders[:i], s.Orders[i+1:]...)
			return s.Save()
		}
	}

	return models.ErrorOrderNotFound
}

func (s *Storage) Update(updatedOrder *Order) error {
	for i, order := range s.Orders {
		if order.ID == updatedOrder.ID {
			s.Orders[i] = updatedOrder
			return s.Save()
		}
	}

	return models.ErrorOrderNotFound
}

func (s *Storage) Add(order *Order) error {
	exists, err := s.GetByID(order.ID)
	if err == nil && exists != nil {
		return models.ErrorOrderAlreadyExists
	}

	s.Orders = append(s.Orders, order)

	return s.Save()
}
