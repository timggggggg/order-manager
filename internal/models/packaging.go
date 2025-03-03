package models

type PackagingType string

const (
	PackagingDefault PackagingType = ""
	PackagingBag     PackagingType = "bag"
	PackagingBox     PackagingType = "box"
	PackagingFilm    PackagingType = "film"
)
