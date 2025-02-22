package repository

import (
	"context"
	"github.com/personal-project/pitch-league/models"

	"github.com/uptrace/bun"
)

type IUserRepository interface {
	IBaseRepository[models.User]
	GetByEmail(ctx context.Context, email string) (models.User, error)
}

type UserRepository struct {
	BaseRepository[models.User]
}

func NewUserRepository(db *bun.DB) UserRepository {
	return UserRepository{
		BaseRepository[models.User]{db: db},
	}
}

func (r UserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.db.NewSelect().
		Model(&user).
		Where("email = ?", email).
		Scan(ctx)
	return user, err
}
