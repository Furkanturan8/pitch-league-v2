package repository

import (
	"context"
	"errors"

	"github.com/personal-project/pitch-league/models"
	"github.com/uptrace/bun"
)

type IGameRepository interface {
	IBaseRepository[models.Game]
	GetAllGame(ctx context.Context) ([]models.Game, error)
	GetByGameID(ctx context.Context, id int64) (*models.Game, error)
	DeleteByGameID(ctx context.Context, id int64) error
	UpdateGame(ctx context.Context, m models.Game) error
	CreateGame(ctx context.Context, game models.Game) error
}

type GameRepository struct {
	BaseRepository[models.Game]
}

func NewGameRepository(db *bun.DB) IGameRepository {
	return &GameRepository{
		BaseRepository: BaseRepository[models.Game]{
			db: db,
		},
	}
}

func (r GameRepository) GetAllGame(ctx context.Context) ([]models.Game, error) {
	var games []models.Game
	err := r.db.NewSelect().
		Model(&games).
		Relation("Host").
		Relation("Field").
		Scan(ctx)
	return games, err
}

func (r GameRepository) GetByGameID(ctx context.Context, id int64) (*models.Game, error) {
	game := new(models.Game)
	err := r.db.NewSelect().
		Model(game).
		Relation("Host").
		Relation("Field").
		Where("g.id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return game, nil
}

func (r GameRepository) DeleteByGameID(ctx context.Context, id int64) error {
	result, err := r.db.NewDelete().
		Model((*models.Game)(nil)).
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
		return errors.New("game not found")
	}

	return nil
}

func (r GameRepository) UpdateGame(ctx context.Context, m models.Game) error {
	// Host kontrolü
	var hostExists bool
	err := r.db.NewSelect().
		Model((*models.User)(nil)).
		Where("id = ?", m.HostID).
		Scan(ctx, &hostExists)
	if err != nil {
		return err
	}
	if !hostExists {
		return errors.New("geçersiz hostID: böyle bir kullanıcı mevcut değil")
	}

	// Field kontrolü
	var fieldExists bool
	err = r.db.NewSelect().
		Model((*models.Field)(nil)).
		Where("id = ?", m.FieldID).
		Scan(ctx, &fieldExists)
	if err != nil {
		return err
	}
	if !fieldExists {
		return errors.New("geçersiz fieldID: böyle bir saha mevcut değil")
	}

	_, err = r.db.NewUpdate().
		Model(&m).
		WherePK().
		Exec(ctx)
	return err
}

func (r GameRepository) CreateGame(ctx context.Context, game models.Game) error {
	_, err := r.db.NewInsert().
		Model(&game).
		Exec(ctx)
	return err
}
