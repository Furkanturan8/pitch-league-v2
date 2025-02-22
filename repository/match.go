package repository

import (
	"context"
	"errors"

	"github.com/personal-project/pitch-league/models"
	"github.com/uptrace/bun"
)

type IMatchRepository interface {
	IBaseRepository[models.Match]
	GetAllMatch(ctx context.Context) ([]models.Match, error)
	GetByMatchID(ctx context.Context, id int64) (*models.Match, error)
	DeleteByMatchID(ctx context.Context, id int64) error
	UpdateMatch(ctx context.Context, m models.Match) error
	CreateMatch(ctx context.Context, match models.Match) error
	UpdateLeagueStandings(ctx context.Context, match models.Match) error
	RecalculateRankings(ctx context.Context, leagueID uint) error
}

type MatchRepository struct {
	BaseRepository[models.Match]
}

func NewMatchRepository(db *bun.DB) IMatchRepository {
	return &MatchRepository{
		BaseRepository: BaseRepository[models.Match]{
			db: db,
		},
	}
}

func (r MatchRepository) GetAllMatch(ctx context.Context) ([]models.Match, error) {
	var matches []models.Match
	err := r.db.NewSelect().
		Model(&matches).
		Relation("League").
		Relation("HomeTeam").
		Relation("AwayTeam").
		Relation("Game").
		Relation("Game.Host").
		Relation("Game.Field").
		Relation("HomeTeam.Captain").
		Relation("AwayTeam.Captain").
		Scan(ctx)
	return matches, err
}

func (r MatchRepository) GetByMatchID(ctx context.Context, id int64) (*models.Match, error) {
	match := new(models.Match)
	err := r.db.NewSelect().
		Model(match).
		Relation("League").
		Relation("HomeTeam").
		Relation("AwayTeam").
		Relation("Game").
		Relation("Game.Host").
		Relation("Game.Field").
		Relation("HomeTeam.Captain").
		Relation("AwayTeam.Captain").
		Where("m.id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return match, nil
}

func (r MatchRepository) DeleteByMatchID(ctx context.Context, id int64) error {
	result, err := r.db.NewDelete().
		Model((*models.Match)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("match not found")
	}

	return nil
}

func (r MatchRepository) UpdateMatch(ctx context.Context, m models.Match) error {
	_, err := r.db.NewUpdate().
		Model(&m).
		WherePK().
		Exec(ctx)
	return err
}

func (r MatchRepository) CreateMatch(ctx context.Context, match models.Match) error {
	_, err := r.db.NewInsert().
		Model(&match).
		Exec(ctx)
	return err
}

func (r MatchRepository) UpdateLeagueStandings(ctx context.Context, match models.Match) error {
	var homeStanding, awayStanding models.LeagueTeam

	// Ev sahibi takımın puan durumu
	err := r.db.NewSelect().
		Model(&homeStanding).
		Where("league_id = ? AND team_id = ?", match.LeagueID, match.HomeTeamID).
		Scan(ctx)

	if err != nil {
		homeStanding = models.LeagueTeam{
			LeagueID: match.LeagueID,
			TeamID:   match.HomeTeamID,
			Points:   0,
			Rank:     0,
		}
	}

	// Deplasman takımının puan durumu
	err = r.db.NewSelect().
		Model(&awayStanding).
		Where("league_id = ? AND team_id = ?", match.LeagueID, match.AwayTeamID).
		Scan(ctx)

	if err != nil {
		awayStanding = models.LeagueTeam{
			LeagueID: match.LeagueID,
			TeamID:   match.AwayTeamID,
			Points:   0,
			Rank:     0,
		}
	}

	// Ev sahibi takımın puanlarını güncelle
	if match.HomeScore > match.AwayScore {
		homeStanding.Points += 3
	} else if match.HomeScore == match.AwayScore {
		homeStanding.Points += 1
	}

	// Deplasman takımının puanlarını güncelle
	if match.AwayScore > match.HomeScore {
		awayStanding.Points += 3
	} else if match.AwayScore == match.HomeScore {
		awayStanding.Points += 1
	}

	// Puan durumlarını kaydet
	_, err = r.db.NewInsert().
		Model(&homeStanding).
		On("CONFLICT (league_id, team_id) DO UPDATE").
		Set("points = EXCLUDED.points").
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = r.db.NewInsert().
		Model(&awayStanding).
		On("CONFLICT (league_id, team_id) DO UPDATE").
		Set("points = EXCLUDED.points").
		Exec(ctx)
	if err != nil {
		return err
	}

	// Sıralamaları yeniden hesapla
	return r.RecalculateRankings(ctx, match.LeagueID)
}

func (r MatchRepository) RecalculateRankings(ctx context.Context, leagueID uint) error {
	var standings []models.LeagueTeam

	// Lig için tüm puan durumlarını getir
	err := r.db.NewSelect().
		Model(&standings).
		Where("league_id = ?", leagueID).
		Order("points DESC").
		Scan(ctx)
	if err != nil {
		return err
	}

	// Sıralı puan durumlarına göre sıralamaları güncelle
	for rank, standing := range standings {
		standing.Rank = int64(rank + 1)
		_, err := r.db.NewUpdate().
			Model(&standing).
			WherePK().
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
