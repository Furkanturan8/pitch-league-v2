package models

import (
	"strconv"

	"github.com/uptrace/bun"
)

type LeagueTeam struct {
	bun.BaseModel `bun:"table:league_teams,alias:lt"`
	ID            int64   `bun:"id,pk,autoincrement" json:"id"`
	LeagueID      uint    `bun:"league_id,notnull" json:"league_id"`
	TeamID        uint    `bun:"team_id,notnull" json:"team_id"`
	Points        int64   `bun:"points,default:0" json:"points"`
	Rank          int64   `bun:"rank" json:"rank"`
	League        *League `bun:"rel:has-one,join:league_id=id" json:"league"`
	Team          *Team   `bun:"rel:has-one,join:team_id=id" json:"team"`
}

type LeagueTeamCreateVM struct {
	LeagueID uint  `json:"league_id" validate:"required"`
	TeamID   uint  `json:"team_id" validate:"required"`
	Points   int64 `json:"points"`
	Rank     int64 `json:"rank"`
}

func (vm LeagueTeamCreateVM) ToDBModel(m LeagueTeam) LeagueTeam {
	m.LeagueID = vm.LeagueID
	m.TeamID = vm.TeamID
	m.Points = vm.Points
	m.Rank = vm.Rank
	return m
}

type LeagueTeamDetailVM struct {
	ID       int64   `json:"id"`
	LeagueID uint    `json:"league_id"`
	TeamID   uint    `json:"team_id"`
	Points   int64   `json:"points"`
	Rank     int64   `json:"rank"`
	League   *League `json:"league"`
	Team     *Team   `json:"team"`
}

func (vm LeagueTeamDetailVM) FromDBModel(m LeagueTeam) LeagueTeamDetailVM {
	vm.ID = m.ID
	vm.LeagueID = m.LeagueID
	vm.TeamID = m.TeamID
	vm.Points = m.Points
	vm.Rank = m.Rank
	vm.League = m.League
	vm.Team = m.Team
	return vm
}

func (LeagueTeam) ModelName() string {
	return "league_teams"
}

func (lt LeagueTeam) String() string {
	if lt.Team != nil && lt.League != nil {
		return lt.Team.Name + " " + lt.League.Name + " " + strconv.Itoa(int(lt.Points)) + " puan ile " + strconv.Itoa(int(lt.Rank)) + " sırada!"
	}
	return "League Team"
}
