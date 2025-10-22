package model

import (
	"time"
)

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Cuntry  string
	Website string
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
