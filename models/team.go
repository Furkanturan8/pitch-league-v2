package models

import (
	"github.com/uptrace/bun"
)

type Team struct {
	bun.BaseModel `bun:"table:teams,alias:t"`
	ID            int64  `bun:"id,pk,autoincrement" json:"id"`
	Name          string `bun:"name,notnull" json:"name"`
	Capacity      int64  `bun:"capacity,notnull" json:"capacity"`
	CaptainID     int64  `bun:"captain_id,notnull" json:"captain_id"`
	Captain       *User  `bun:"rel:has-one,join:captain_id=id" json:"captain"`
}

type TeamCreateVM struct {
	Name      string `json:"name" validate:"required,max=100"`
	Capacity  int64  `json:"capacity" validate:"required,max=100"`
	CaptainID int64  `json:"captain_id" validate:"required"`
}

func (vm TeamCreateVM) ToDBModel(m Team) Team {
	m.Name = vm.Name
	m.Capacity = vm.Capacity
	m.CaptainID = vm.CaptainID
	return m
}

type TeamDetailVM struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	Capacity  int64       `json:"capacity"`
	CaptainID int64       `json:"captain_id"`
	Captain   interface{} `json:"captain"`
}

func (vm TeamDetailVM) FromDBModel(m Team) TeamDetailVM {
	vm.ID = m.ID
	vm.Name = m.Name
	vm.Capacity = m.Capacity
	vm.CaptainID = m.CaptainID
	vm.Captain = m.Captain
	return vm
}

func (Team) ModelName() string {
	return "teams"
}

func (t Team) String() string {
	if t.Captain != nil {
		return "Takım: " + t.Name + " Kaptan: " + t.Captain.Name
	}
	return "Takım: " + t.Name
}
