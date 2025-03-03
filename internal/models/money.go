package models

import (
	"fmt"
	"strconv"
	"strings"
)

type Money struct {
	Amount int64
}

func NewMoney(str string) (*Money, error) {
	if strings.Count(str, ".") == 0 {
		str += ".0"
	}

	temp := strings.Split(str, ".")

	if len(temp) != 2 {
		return nil, ErrorInvalidAmountOfMoney
	}

	f, err := strconv.ParseInt(temp[0], 10, 64)
	if err != nil {
		return nil, ErrorInvalidAmountOfMoney
	}

	if len(temp[1]) > 2 {
		return nil, ErrorInvalidAmountOfMoney
	}

	s, err := strconv.ParseInt(temp[1], 10, 64)
	if err != nil {
		return nil, ErrorInvalidAmountOfMoney
	}

	return &Money{f*100 + s}, nil
}

func (m *Money) Add(otherAmount int64) {
	m.Amount += otherAmount
}

func (m *Money) String() string {
	return fmt.Sprintf("%d.%d", m.Amount/100, m.Amount%100)
}
