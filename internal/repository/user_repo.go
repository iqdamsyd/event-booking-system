package repository

import (
	"context"
	"event-booking-system/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
	`, user.Name, user.Email, user.Password)

	return err
}

func (r *UserRepository) List(ctx context.Context) ([]models.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, email, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user models.User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.RequestUpdateUser) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET name = $1, email = $2, password = $3, updated_at = NOW()
		WHERE id = $4
	`, user.Name, user.Email, user.Password, user.ID)

	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM users WHERE id = $1
	`, id)

	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, email, password, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user models.User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
	}

	return &user, nil
}
