package models

type PackagingType string

const (
	PackagingDefault PackagingType = "default"
	PackagingBag     PackagingType = "bag"
	PackagingBox     PackagingType = "box"
	PackagingFilm    PackagingType = "film"
)
