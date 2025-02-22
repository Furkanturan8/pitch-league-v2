package models

import (
	"time"

	"github.com/uptrace/bun"
)

type League struct {
	bun.BaseModel `bun:"table:leagues,alias:l"`
	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	Name          string    `bun:"name,notnull,unique" json:"name"`
	Location      string    `bun:"location,notnull" json:"location"`
	StartDate     time.Time `bun:"start_date,notnull" json:"start_date"`
	EndDate       time.Time `bun:"end_date,notnull" json:"end_date"`
}

type LeagueCreateVM struct {
	Name      string    `json:"name" validate:"required"`
	Location  string    `json:"location" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
}

func (vm LeagueCreateVM) ToDBModel(m League) League {
	m.Name = vm.Name
	m.Location = vm.Location
	m.StartDate = vm.StartDate
	m.EndDate = vm.EndDate
	return m
}

type LeagueDetailVM struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func (vm LeagueDetailVM) FromDBModel(m League) LeagueDetailVM {
	vm.ID = m.ID
	vm.Name = m.Name
	vm.Location = m.Location
	vm.StartDate = m.StartDate
	vm.EndDate = m.EndDate
	return vm
}

func (League) ModelName() string {
	return "leagues"
}

func (l League) String() string {
	return l.Name + " " + l.Location
}
