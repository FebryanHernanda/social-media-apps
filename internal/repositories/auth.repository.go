package repositories

import (
	"context"

	"github.com/febryanhernanda/social-media-apps/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	DB *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		DB: db,
	}
}

func (r *AuthRepository) RegisterUser(ctx context.Context, req *models.RegisterUser) (*models.User, error) {
	query := `
		INSERT INTO users (email, password, name)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	var user models.User
	err := r.DB.QueryRow(ctx, query, req.Email, req.Password, req.Name).
		Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) LoginUser(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password FROM users WHERE email=$1`

	err := r.DB.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
