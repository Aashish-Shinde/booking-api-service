package repository

import (
	"context"
	"database/sql"

	"github.com/booking-api-service/internal/model"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) (int64, error) {
	query := `INSERT INTO users (name, timezone) VALUES (?, ?)`
	result, err := r.db.ExecContext(ctx, query, user.Name, user.Timezone)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `SELECT id, name, timezone, created_at, updated_at, deleted FROM users WHERE id = ? AND deleted = FALSE`
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Timezone, &user.CreatedAt, &user.UpdatedAt, &user.Deleted,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]*model.User, error) {
	query := `SELECT id, name, timezone, created_at, updated_at, deleted FROM users WHERE deleted = FALSE`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Timezone, &user.CreatedAt, &user.UpdatedAt, &user.Deleted)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `UPDATE users SET name = ?, timezone = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, user.Name, user.Timezone, user.ID)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE users SET deleted = TRUE WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
