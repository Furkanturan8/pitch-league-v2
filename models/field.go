package models

import (
	"github.com/uptrace/bun"
)

type Field struct {
	bun.BaseModel `bun:"table:fields,alias:f"`
	ID            int64   `bun:"id,pk,autoincrement" json:"id"`
	Name          string  `bun:"name,notnull" json:"name"`
	Location      string  `bun:"location,notnull" json:"location"`
	PricePerHour  float64 `bun:"price_per_hour,notnull" json:"price_per_hour"`
	Capacity      int64   `bun:"capacity,notnull" json:"capacity"`
	Available     bool    `bun:"available,notnull" json:"available"`
}

type FieldCreateVM struct {
	Name         string  `json:"name" validate:"required,max=100"`
	Location     string  `json:"location" validate:"required,max=255"`
	PricePerHour float64 `json:"price_per_hour" validate:"required"`
	Capacity     int64   `json:"capacity" validate:"required"`
	Available    bool    `json:"available" validate:"omitempty"`
}

func (vm FieldCreateVM) ToDBModel(m Field) Field {
	m.Name = vm.Name
	m.Location = vm.Location
	m.PricePerHour = vm.PricePerHour
	m.Capacity = vm.Capacity
	m.Available = vm.Available
	return m
}

type FieldDetailVM struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Location     string  `json:"location"`
	PricePerHour float64 `json:"price_per_hour"`
	Capacity     int64   `json:"capacity"`
	Available    bool    `json:"available"`
}

func (vm FieldDetailVM) FromDBModel(m Field) FieldDetailVM {
	vm.ID = m.ID
	vm.Name = m.Name
	vm.Location = m.Location
	vm.PricePerHour = m.PricePerHour
	vm.Capacity = m.Capacity
	vm.Available = m.Available
	return vm
}

func (Field) ModelName() string {
	return "fields"
}

func (f Field) String() string {
	return f.Name + " " + f.Location
}
