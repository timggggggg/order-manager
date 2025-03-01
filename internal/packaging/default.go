package packaging

import (
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

const (
	DefaultCost = 0.0
)

type DefaultPackaging struct{}

func (b *DefaultPackaging) CalculateCost(weight float64) (float64, error) {
	return DefaultCost, nil
}

func (b *DefaultPackaging) Type() models.PackagingType {
	return models.PackagingDefault
}
