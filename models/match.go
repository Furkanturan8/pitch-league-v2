package models

import (
	"time"

	"github.com/uptrace/bun"
)

type MatchStatus string

const (
	MatchStatusScheduled MatchStatus = "SCHEDULED"
	MatchStatusCompleted MatchStatus = "COMPLETED"
)

type Match struct {
	bun.BaseModel `bun:"table:matches,alias:m"`
	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	LeagueID      uint      `bun:"league_id,notnull" json:"league_id"`
	HomeTeamID    uint      `bun:"home_team_id,notnull" json:"home_team_id"`
	AwayTeamID    uint      `bun:"away_team_id,notnull" json:"away_team_id"`
	MatchTime     time.Time `bun:"match_time,notnull" json:"match_time"`
	HomeScore     int64     `bun:"home_score,default:0" json:"home_score"`
	AwayScore     int64     `bun:"away_score,default:0" json:"away_score"`
	Status        string    `bun:"status,notnull" json:"status"`
	GameID        uint      `bun:"game_id,notnull" json:"game_id"`
	Game          *Game     `bun:"rel:has-one,join:game_id=id" json:"game"`
	League        *League   `bun:"rel:has-one,join:league_id=id" json:"league"`
	HomeTeam      *Team     `bun:"rel:has-one,join:home_team_id=id" json:"home_team"`
	AwayTeam      *Team     `bun:"rel:has-one,join:away_team_id=id" json:"away_team"`
}

type MatchCreateVM struct {
	LeagueID   uint      `json:"league_id" validate:"required"`
	HomeTeamID uint      `json:"home_team_id" validate:"required"`
	AwayTeamID uint      `json:"away_team_id" validate:"required"`
	MatchTime  time.Time `json:"match_time" validate:"required"`
	GameID     uint      `json:"game_id" validate:"required"`
	HomeScore  int64     `json:"home_score"`
	AwayScore  int64     `json:"away_score"`
	Status     string    `json:"status"`
}

func (vm MatchCreateVM) ToDBModel(m Match) Match {
	m.LeagueID = vm.LeagueID
	m.HomeTeamID = vm.HomeTeamID
	m.AwayTeamID = vm.AwayTeamID
	m.MatchTime = vm.MatchTime
	m.GameID = vm.GameID
	m.HomeScore = vm.HomeScore
	m.AwayScore = vm.AwayScore
	m.Status = vm.Status
	return m
}

type MatchDetailVM struct {
	ID         int64     `json:"id"`
	LeagueID   uint      `json:"league_id"`
	HomeTeamID uint      `json:"home_team_id"`
	AwayTeamID uint      `json:"away_team_id"`
	MatchTime  time.Time `json:"match_time"`
	HomeScore  int64     `json:"home_score"`
	AwayScore  int64     `json:"away_score"`
	Status     string    `json:"status"`
	GameID     uint      `json:"game_id"`
	Game       *Game     `json:"game"`
	League     *League   `json:"league"`
	HomeTeam   *Team     `json:"home_team"`
	AwayTeam   *Team     `json:"away_team"`
}

func (vm MatchDetailVM) FromDBModel(m Match) MatchDetailVM {
	vm.ID = m.ID
	vm.LeagueID = m.LeagueID
	vm.HomeTeamID = m.HomeTeamID
	vm.AwayTeamID = m.AwayTeamID
	vm.MatchTime = m.MatchTime
	vm.HomeScore = m.HomeScore
	vm.AwayScore = m.AwayScore
	vm.Status = m.Status
	vm.GameID = m.GameID
	vm.Game = m.Game
	vm.League = m.League
	vm.HomeTeam = m.HomeTeam
	vm.AwayTeam = m.AwayTeam
	return vm
}

func (Match) ModelName() string {
	return "matches"
}

func (m Match) String() string {
	if m.HomeTeam != nil && m.AwayTeam != nil {
		return m.HomeTeam.Name + " vs " + m.AwayTeam.Name
	}
	return "Match"
}
