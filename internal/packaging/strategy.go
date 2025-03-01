package packaging

import (
	"fmt"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type Strategy interface {
	CalculateCost(weight float64) (float64, error)
	Type() models.PackagingType
}

type PackagingStrategyMap map[string]func() Strategy

var PackagingStrategies = PackagingStrategyMap{
	string(models.PackagingBag):  func() Strategy { return &BagPackaging{} },
	string(models.PackagingBox):  func() Strategy { return &BoxPackaging{} },
	string(models.PackagingFilm): func() Strategy { return &FilmPackaging{} },
}

var ExtraPackagingStrategies = PackagingStrategyMap{
	string(models.PackagingDefault): func() Strategy { return &DefaultPackaging{} },
	string(models.PackagingFilm):    func() Strategy { return &FilmPackaging{} },
}

func NewPackagingStrategy(packaging string, strategies PackagingStrategyMap) (Strategy, error) {
	createStrategy, exists := strategies[packaging]
	if !exists {
		return nil, fmt.Errorf("invalid package type: %s", packaging)
	}
	return createStrategy(), nil
}
