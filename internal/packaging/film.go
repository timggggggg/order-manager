package packaging

import (
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

const (
	FilmCost = 1.0
)

type FilmPackaging struct{}

func (b *FilmPackaging) CalculateCost(weight float64) (float64, error) {
	return FilmCost, nil
}

func (b *FilmPackaging) Type() models.PackagingType {
	return models.PackagingFilm
}
