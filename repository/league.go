package repository

import (
	"context"
	"errors"

	"github.com/personal-project/pitch-league/models"
	"github.com/uptrace/bun"
)

type ILeagueRepository interface {
	IBaseRepository[models.League]
	GetAllLeague(ctx context.Context) ([]models.League, error)
	GetByLeagueID(ctx context.Context, id int64) (*models.League, error)
	DeleteByLeagueID(ctx context.Context, id int64) error
	UpdateLeague(ctx context.Context, m models.League) error
	CreateLeague(ctx context.Context, league models.League) error
}

type LeagueRepository struct {
	BaseRepository[models.League]
}

func NewLeagueRepository(db *bun.DB) ILeagueRepository {
	return &LeagueRepository{
		BaseRepository: BaseRepository[models.League]{
			db: db,
		},
	}
}

func (r LeagueRepository) GetAllLeague(ctx context.Context) ([]models.League, error) {
	var leagues []models.League
	err := r.db.NewSelect().
		Model(&leagues).
		Scan(ctx)
	return leagues, err
}

func (r LeagueRepository) GetByLeagueID(ctx context.Context, id int64) (*models.League, error) {
	league := new(models.League)
	err := r.db.NewSelect().
		Model(league).
		Where("l.id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return league, nil
}

func (r LeagueRepository) DeleteByLeagueID(ctx context.Context, id int64) error {
	result, err := r.db.NewDelete().
		Model((*models.League)(nil)).
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
		return errors.New("league not found")
	}

	return nil
}

func (r LeagueRepository) UpdateLeague(ctx context.Context, m models.League) error {
	_, err := r.db.NewUpdate().
		Model(&m).
		WherePK().
		Exec(ctx)
	return err
}

func (r LeagueRepository) CreateLeague(ctx context.Context, league models.League) error {
	_, err := r.db.NewInsert().
		Model(&league).
		Exec(ctx)
	return err
}
