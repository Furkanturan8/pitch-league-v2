package models

import (
	"time"

	"github.com/uptrace/bun"
)

type GameStatus string

const (
	GameStatusPending   GameStatus = "PENDING"
	GameStatusAccepted  GameStatus = "ACCEPTED"
	GameStatusRejected  GameStatus = "REJECTED"
	GameStatusCancelled GameStatus = "CANCELLED"
	GameStatusFinished  GameStatus = "FINISHED"
)

type Game struct {
	bun.BaseModel `bun:"table:games,alias:g"`
	ID            int64      `bun:"id,pk,autoincrement" json:"id"`
	FieldID       uint       `bun:"field_id,notnull" json:"field_id"`
	HostID        uint       `bun:"host_id,notnull" json:"host_id"`
	StartTime     time.Time  `bun:"start_time,notnull" json:"start_time"`
	EndTime       time.Time  `bun:"end_time,notnull" json:"end_time"`
	MaxPlayers    int64      `bun:"max_players,notnull" json:"max_players"`
	Status        GameStatus `bun:"status,notnull" json:"status"`
	Host          *User      `bun:"rel:has-one,join:host_id=id" json:"host"`
	Field         *Field     `bun:"rel:has-one,join:field_id=id" json:"field"`
}

type GameCreateVM struct {
	FieldID    uint       `json:"field_id" validate:"required"`
	HostID     uint       `json:"host_id" validate:"required"`
	StartTime  time.Time  `json:"start_time" validate:"required"`
	EndTime    time.Time  `json:"end_time" validate:"required"`
	MaxPlayers int64      `json:"max_players" validate:"required"`
	Status     GameStatus `json:"status" validate:"required"`
}

func (vm GameCreateVM) ToDBModel(m Game) Game {
	m.FieldID = vm.FieldID
	m.HostID = vm.HostID
	m.StartTime = vm.StartTime
	m.EndTime = vm.EndTime
	m.MaxPlayers = vm.MaxPlayers
	m.Status = vm.Status
	return m
}

type GameDetailVM struct {
	ID         int64       `json:"id"`
	FieldID    uint        `json:"field_id"`
	HostID     uint        `json:"host_id"`
	StartTime  time.Time   `json:"start_time"`
	EndTime    time.Time   `json:"end_time"`
	MaxPlayers int64       `json:"max_players"`
	Status     GameStatus  `json:"status"`
	Host       interface{} `json:"host"`
	Field      *Field      `json:"field"`
}

func (vm GameDetailVM) FromDBModel(m Game) GameDetailVM {
	vm.ID = m.ID
	vm.FieldID = m.FieldID
	vm.HostID = m.HostID
	vm.StartTime = m.StartTime
	vm.EndTime = m.EndTime
	vm.MaxPlayers = m.MaxPlayers
	vm.Status = m.Status
	vm.Host = m.Host
	vm.Field = m.Field
	return vm
}

func (Game) ModelName() string {
	return "games"
}

func (g GameStatus) String() string {
	switch g {
	case GameStatusPending:
		return "Pending"
	case GameStatusAccepted:
		return "Accepted"
	case GameStatusRejected:
		return "Rejected"
	case GameStatusCancelled:
		return "Cancelled"
	case GameStatusFinished:
		return "Finished"
	default:
		return "Unknown"
	}
}
