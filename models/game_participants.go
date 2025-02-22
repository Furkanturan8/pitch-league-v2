package models

import (
	"github.com/uptrace/bun"
)

type GameParticipants struct {
	bun.BaseModel `bun:"table:game_participants,alias:gp"`
	ID            int64 `bun:"id,pk,autoincrement" json:"id"`
	GameID        uint  `bun:"game_id,notnull" json:"game_id"`
	UserID        uint  `bun:"user_id,notnull" json:"user_id"`
	TeamID        uint  `bun:"team_id,notnull" json:"team_id"`
	Game          *Game `bun:"rel:has-one,join:game_id=id" json:"game"`
	User          *User `bun:"rel:has-one,join:user_id=id" json:"user"`
	Team          *Team `bun:"rel:has-one,join:team_id=id" json:"team"`
}

type GameParticipantsCreateVM struct {
	GameID uint `json:"game_id" validate:"required"`
	UserID uint `json:"user_id" validate:"required"`
	TeamID uint `json:"team_id" validate:"required"`
}

func (vm GameParticipantsCreateVM) ToDBModel(m GameParticipants) GameParticipants {
	m.GameID = vm.GameID
	m.UserID = vm.UserID
	m.TeamID = vm.TeamID
	return m
}

type GameParticipantsDetailVM struct {
	ID     int64 `json:"id"`
	GameID uint  `json:"game_id"`
	UserID uint  `json:"user_id"`
	TeamID uint  `json:"team_id"`
	Game   *Game `json:"game"`
	Team   *Team `json:"team"`
	User   *User `json:"user"`
}

func (vm GameParticipantsDetailVM) FromDBModel(m GameParticipants) GameParticipantsDetailVM {
	vm.ID = m.ID
	vm.GameID = m.GameID
	vm.UserID = m.UserID
	vm.TeamID = m.TeamID
	vm.Game = m.Game
	vm.Team = m.Team
	vm.User = m.User
	return vm
}

type GameParticipantsUsersVM struct {
	GameID uint   `json:"game_id"`
	Users  []User `json:"users"`
}

func (vm GameParticipantsUsersVM) FromDBModel(gameID uint, users []User) GameParticipantsUsersVM {
	vm.GameID = gameID
	vm.Users = users
	return vm
}

func (GameParticipants) ModelName() string {
	return "game_participants"
}
