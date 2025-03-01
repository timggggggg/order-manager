package packaging

import (
	"fmt"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

const (
	BoxCost      = 20.0
	BoxMaxWeight = 30.0
)

var ErrorBoxPackaging = fmt.Errorf("box is only available for orders under %.2f kg", BoxMaxWeight)

type BoxPackaging struct{}

func (b *BoxPackaging) CalculateCost(weight float64) (float64, error) {
	if weight >= BoxMaxWeight {
		return 0, ErrorBoxPackaging
	}

	return BoxCost, nil
}

func (b *BoxPackaging) Type() models.PackagingType {
	return models.PackagingBox
}
