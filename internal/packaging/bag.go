package packaging

import (
	"fmt"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

const (
	BagCost      = 5
	BagMaxWeight = 10.0
)

var ErrorBagPackaging = fmt.Errorf("bag is only available for orders under %.2f kg", BagMaxWeight)

type BagPackaging struct{}

func (b *BagPackaging) CalculateCost(weight float64) (*models.Money, error) {
	if weight >= BagMaxWeight {
		return &models.Money{Amount: 0}, ErrorBagPackaging
	}

	return &models.Money{Amount: BagCost}, nil
}

func (b *BagPackaging) Type() models.PackagingType {
	return models.PackagingBag
}
