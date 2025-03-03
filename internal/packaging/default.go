package packaging

import (
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

const (
	DefaultCost = 0
)

type DefaultPackaging struct{}

func (b *DefaultPackaging) CalculateCost(weight float64) (*models.Money, error) {
	return &models.Money{Amount: DefaultCost}, nil
}

func (b *DefaultPackaging) Type() models.PackagingType {
	return models.PackagingDefault
}
