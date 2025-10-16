package model

import "time"

type Dimensions struct {
	length int32
	width  int32
	height int32
	weight int32
}

type Manufacturer struct {
	name    string
	cuntry  string
	website string
}

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      int32
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Filter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []int32
	ManufacturerCountries []string
	Tags                  []string
}
