package repository

import (
	"context"
	"errors"

	"github.com/personal-project/pitch-league/models"
	"github.com/uptrace/bun"
)

type IFieldRepository interface {
	IBaseRepository[models.Field]
	GetAllField(ctx context.Context) ([]models.Field, error)
	GetByFieldID(ctx context.Context, id int64) (*models.Field, error)
	DeleteByFieldID(ctx context.Context, id int64) error
	UpdateField(ctx context.Context, m models.Field) error
	CreateField(ctx context.Context, field models.Field) error
}

type FieldRepository struct {
	BaseRepository[models.Field]
}

func NewFieldRepository(db *bun.DB) IFieldRepository {
	return &FieldRepository{
		BaseRepository: BaseRepository[models.Field]{
			db: db,
		},
	}
}

func (r FieldRepository) GetAllField(ctx context.Context) ([]models.Field, error) {
	var fields []models.Field
	err := r.db.NewSelect().
		Model(&fields).
		Scan(ctx)
	return fields, err
}

func (r FieldRepository) GetByFieldID(ctx context.Context, id int64) (*models.Field, error) {
	field := new(models.Field)
	err := r.db.NewSelect().
		Model(field).
		Where("f.id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return field, nil
}

func (r FieldRepository) DeleteByFieldID(ctx context.Context, id int64) error {
	result, err := r.db.NewDelete().
		Model((*models.Field)(nil)).
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
		return errors.New("field not found")
	}

	return nil
}

func (r FieldRepository) UpdateField(ctx context.Context, m models.Field) error {
	_, err := r.db.NewUpdate().
		Model(&m).
		WherePK().
		Exec(ctx)
	return err
}

func (r FieldRepository) CreateField(ctx context.Context, field models.Field) error {
	_, err := r.db.NewInsert().
		Model(&field).
		Exec(ctx)
	return err
}
