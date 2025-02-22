package repository

import (
	"context"
	"errors"

	"github.com/personal-project/pitch-league/models"
	"github.com/uptrace/bun"
)

type ITeamRepository interface {
	IBaseRepository[models.Team]
	AddUserToTeam(ctx context.Context, userID, teamID int64) error
	GetAllTeam(ctx context.Context) ([]models.Team, error)
	GetByTeamID(ctx context.Context, id int64) (*models.Team, error)
	DeleteByTeamID(ctx context.Context, id int64) error
	UpdateTeam(ctx context.Context, m models.Team) error
	CreateTeam(ctx context.Context, team models.Team) error
}

type TeamRepository struct {
	BaseRepository[models.Team]
}

func NewTeamRepository(db *bun.DB) ITeamRepository {
	return &TeamRepository{
		BaseRepository: BaseRepository[models.Team]{
			db: db,
		},
	}
}

func (r TeamRepository) GetAllTeam(ctx context.Context) ([]models.Team, error) {
	var teams []models.Team
	err := r.db.NewSelect().
		Model(&teams).
		Relation("Captain").
		Scan(ctx)
	return teams, err
}

func (r TeamRepository) GetByTeamID(ctx context.Context, id int64) (*models.Team, error) {
	team := new(models.Team)
	err := r.db.NewSelect().
		Model(team).
		Relation("Captain").
		Where("t.id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return team, nil
}

func (r TeamRepository) DeleteByTeamID(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().
		Model((*models.Team)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r TeamRepository) UpdateTeam(ctx context.Context, m models.Team) error {
	_, err := r.db.NewUpdate().
		Model(&m).
		WherePK().
		Exec(ctx)
	return err
}

func (r TeamRepository) AddUserToTeam(ctx context.Context, userID, teamID int64) error {
	// Takımı kontrol et
	team := new(models.Team)
	err := r.db.NewSelect().
		Model(team).
		Where("t.id = ?", teamID).
		Scan(ctx)

	if err != nil {
		return err
	}

	if team.Capacity <= 0 {
		return errors.New("team capacity is full")
	}

	// Kullanıcıyı güncelle
	_, err = r.db.NewUpdate().
		Model((*models.User)(nil)).
		Set("team_id = ?", teamID).
		Where("id = ?", userID).
		Exec(ctx)

	if err != nil {
		return err
	}

	// Takım kapasitesini güncelle
	_, err = r.db.NewUpdate().
		Model(team).
		Set("capacity = capacity - 1").
		Where("t.id = ?", teamID).
		Exec(ctx)

	return err
}

func (r TeamRepository) CreateTeam(ctx context.Context, team models.Team) error {
	// Önce takımı oluştur
	_, err := r.db.NewInsert().
		Model(&team).
		Exec(ctx)
	if err != nil {
		return err
	}

	// Oluşturulan takımı Captain ilişkisiyle birlikte yükle
	err = r.db.NewSelect().
		Model(&team).
		Relation("Captain").
		Where("t.id = ?", team.ID).
		Scan(ctx)

	return err
}
