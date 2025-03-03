package packaging

import (
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

const (
	FilmCost = 1
)

type FilmPackaging struct{}

func (b *FilmPackaging) CalculateCost(weight float64) (*models.Money, error) {
	return &models.Money{Amount: FilmCost}, nil
}

func (b *FilmPackaging) Type() models.PackagingType {
	return models.PackagingFilm
}
