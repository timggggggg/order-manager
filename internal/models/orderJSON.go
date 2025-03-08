package models

type OrderJSON struct {
	ID                  int64   `json:"id"`
	UserID              int64   `json:"user_id"`
	StorageDurationDays int64   `json:"storage_duration"`
	Weight              float64 `json:"weight"`
	Cost                string  `json:"cost"`
	Package             string  `json:"package"`
	ExtraPackage        string  `json:"extra_package,omitempty"`
}
