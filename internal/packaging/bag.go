package packaging

import (
	"fmt"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

const (
	BagCost      = 5.0
	BagMaxWeight = 10.0
)

var ErrorBagPackaging = fmt.Errorf("bag is only available for orders under %.2f kg", BagMaxWeight)

type BagPackaging struct{}

func (b *BagPackaging) CalculateCost(weight float64) (float64, error) {
	if weight >= BagMaxWeight {
		return 0, ErrorBagPackaging
	}

	return BagCost, nil
}

func (b *BagPackaging) Type() models.PackagingType {
	return models.PackagingBag
}
