package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/personal-project/pitch-league/models"
	"github.com/uptrace/bun"
)

type IGameParticipantsRepository interface {
	IBaseRepository[models.GameParticipants]
	GetAllGameParticipants(ctx context.Context) ([]models.GameParticipants, error)
	GetByGameParticipantsID(ctx context.Context, userID int64) (*models.GameParticipants, error)
	GetGameParticipantsUsers(ctx context.Context, gameID uint) ([]models.User, error)
	DeleteByGameParticipantsID(ctx context.Context, id int64) error
	UpdateGameParticipants(ctx context.Context, m models.GameParticipants) error
	CreateGameParticipants(ctx context.Context, gamePart models.GameParticipants) error
	FixGameParticipantsOnTeamChange(ctx context.Context, userID, teamID int64) error
}

type GameParticipantsRepository struct {
	BaseRepository[models.GameParticipants]
}

func NewGameParticipantsRepository(db *bun.DB) IGameParticipantsRepository {
	return &GameParticipantsRepository{
		BaseRepository: BaseRepository[models.GameParticipants]{
			db: db,
		},
	}
}

func (r GameParticipantsRepository) GetAllGameParticipants(ctx context.Context) ([]models.GameParticipants, error) {
	var gameParts []models.GameParticipants
	err := r.db.NewSelect().
		Model(&gameParts).
		Relation("User").
		Relation("Game").
		Relation("Game.Field").
		Relation("Game.Host").
		Relation("Team").
		Relation("Team.Captain").
		Scan(ctx)
	return gameParts, err
}

func (r GameParticipantsRepository) GetByGameParticipantsID(ctx context.Context, userID int64) (*models.GameParticipants, error) {
	gamePart := new(models.GameParticipants)
	err := r.db.NewSelect().
		Model(gamePart).
		Relation("User").
		Relation("Game").
		Relation("Game.Field").
		Relation("Game.Host").
		Relation("Team").
		Relation("Team.Captain").
		Where("gp.user_id = ?", userID).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return gamePart, nil
}

func (r GameParticipantsRepository) GetGameParticipantsUsers(ctx context.Context, gameID uint) ([]models.User, error) {
	var users []models.User
	err := r.db.NewSelect().
		Model(&users).
		Join("JOIN game_participants gp ON gp.user_id = \"user\".id").
		Where("gp.game_id = ?", gameID).
		Scan(ctx)

	return users, err
}

func (r GameParticipantsRepository) DeleteByGameParticipantsID(ctx context.Context, id int64) error {
	result, err := r.db.NewDelete().
		Model((*models.GameParticipants)(nil)).
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
		return errors.New("game participant not found")
	}

	return nil
}

func (r GameParticipantsRepository) UpdateGameParticipants(ctx context.Context, m models.GameParticipants) error {
	_, err := r.db.NewUpdate().
		Model(&m).
		WherePK().
		Exec(ctx)
	return err
}

func (r GameParticipantsRepository) CreateGameParticipants(ctx context.Context, gamePart models.GameParticipants) error {
	// Oyun kontrolü
	var game models.Game
	err := r.db.NewSelect().
		Model(&game).
		Where("id = ?", gamePart.GameID).
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("oyun bulunamadı: %w", err)
	}

	// Kullanıcı kontrolü
	var user models.User
	err = r.db.NewSelect().
		Model(&user).
		Where("id = ?", gamePart.UserID).
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("kullanıcı bulunamadı: %w", err)
	}

	// Takım kontrolü
	var team models.Team
	err = r.db.NewSelect().
		Model(&team).
		Where("id = ?", gamePart.TeamID).
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("takım bulunamadı: %w", err)
	}

	_, err = r.db.NewInsert().
		Model(&gamePart).
		Exec(ctx)
	return err
}

func (r GameParticipantsRepository) FixGameParticipantsOnTeamChange(ctx context.Context, userID, teamID int64) error {
	var gameParts []models.GameParticipants
	err := r.db.NewSelect().
		Model(&gameParts).
		Relation("Game").
		Where("user_id = ? AND team_id = ?", userID, teamID).
		Scan(ctx)
	if err != nil {
		return err
	}

	for _, gamePart := range gameParts {
		if gamePart.Game.Status == models.GameStatusPending {
			_, err = r.db.NewDelete().
				Model((*models.GameParticipants)(nil)).
				Where("game_id = ? AND user_id = ? AND team_id = ?", gamePart.GameID, userID, teamID).
				Exec(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
