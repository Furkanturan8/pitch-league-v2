package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/personal-project/pitch-league/models"
	"github.com/uptrace/bun"
)

type ILeagueTeamRepository interface {
	IBaseRepository[models.LeagueTeam]
	GetAllLeagueTeam(ctx context.Context) ([]models.LeagueTeam, error)
	GetByLeagueTeamID(ctx context.Context, id int64) (*models.LeagueTeam, error)
	GetByLeagueID(ctx context.Context, id int64) ([]models.LeagueTeam, error)
	DeleteByLeagueTeamID(ctx context.Context, id int64) error
	UpdateLeagueTeam(ctx context.Context, m models.LeagueTeam) error
	CreateLeagueTeam(ctx context.Context, leagueTeam models.LeagueTeam) error
}

type LeagueTeamRepository struct {
	BaseRepository[models.LeagueTeam]
}

func NewLeagueTeamRepository(db *bun.DB) ILeagueTeamRepository {
	return &LeagueTeamRepository{
		BaseRepository: BaseRepository[models.LeagueTeam]{
			db: db,
		},
	}
}

func (r LeagueTeamRepository) GetAllLeagueTeam(ctx context.Context) ([]models.LeagueTeam, error) {
	var leagueTeams []models.LeagueTeam
	err := r.db.NewSelect().
		Model(&leagueTeams).
		Relation("Team").
		Relation("League").
		Relation("Team.Captain").
		Order("points DESC").
		Order("rank ASC").
		Scan(ctx)
	return leagueTeams, err
}

func (r LeagueTeamRepository) GetByLeagueTeamID(ctx context.Context, id int64) (*models.LeagueTeam, error) {
	leagueTeam := new(models.LeagueTeam)
	err := r.db.NewSelect().
		Model(leagueTeam).
		Relation("Team").
		Relation("League").
		Relation("Team.Captain").
		Where("lt.team_id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return leagueTeam, nil
}

func (r LeagueTeamRepository) GetByLeagueID(ctx context.Context, id int64) ([]models.LeagueTeam, error) {
	var leagueTeams []models.LeagueTeam
	err := r.db.NewSelect().
		Model(&leagueTeams).
		Relation("Team").
		Relation("League").
		Relation("Team.Captain").
		Where("lt.league_id = ?", id).
		Order("points DESC").
		Order("rank ASC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return leagueTeams, nil
}

func (r LeagueTeamRepository) DeleteByLeagueTeamID(ctx context.Context, id int64) error {
	result, err := r.db.NewDelete().
		Model((*models.LeagueTeam)(nil)).
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
		return errors.New("league team not found")
	}

	return nil
}

func (r LeagueTeamRepository) UpdateLeagueTeam(ctx context.Context, m models.LeagueTeam) error {
	_, err := r.db.NewUpdate().
		Model(&m).
		WherePK().
		Exec(ctx)
	return err
}

func (r LeagueTeamRepository) CreateLeagueTeam(ctx context.Context, leagueTeam models.LeagueTeam) error {
	// Lig kontrolü
	var league models.League
	err := r.db.NewSelect().
		Model(&league).
		Where("id = ?", leagueTeam.LeagueID).
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("lig bulunamadı: %w", err)
	}

	// Takım kontrolü
	var team models.Team
	err = r.db.NewSelect().
		Model(&team).
		Where("id = ?", leagueTeam.TeamID).
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("takım bulunamadı: %w", err)
	}

	_, err = r.db.NewInsert().
		Model(&leagueTeam).
		Exec(ctx)
	return err
}
