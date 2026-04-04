package repository

import (
	"context"
	"database/sql"

	"github.com/booking-api-service/internal/model"
)

type coachRepository struct {
	db *sql.DB
}

func NewCoachRepository(db *sql.DB) CoachRepository {
	return &coachRepository{db: db}
}

func (r *coachRepository) Create(ctx context.Context, coach *model.Coach) (int64, error) {
	query := `INSERT INTO coaches (name, timezone) VALUES (?, ?)`
	result, err := r.db.ExecContext(ctx, query, coach.Name, coach.Timezone)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *coachRepository) GetByID(ctx context.Context, id int64) (*model.Coach, error) {
	query := `SELECT id, name, timezone, created_at, updated_at, deleted FROM coaches WHERE id = ? AND deleted = FALSE`
	coach := &model.Coach{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&coach.ID, &coach.Name, &coach.Timezone, &coach.CreatedAt, &coach.UpdatedAt, &coach.Deleted,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return coach, nil
}

func (r *coachRepository) GetAll(ctx context.Context) ([]*model.Coach, error) {
	query := `SELECT id, name, timezone, created_at, updated_at, deleted FROM coaches WHERE deleted = FALSE`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coaches []*model.Coach
	for rows.Next() {
		coach := &model.Coach{}
		err := rows.Scan(&coach.ID, &coach.Name, &coach.Timezone, &coach.CreatedAt, &coach.UpdatedAt, &coach.Deleted)
		if err != nil {
			return nil, err
		}
		coaches = append(coaches, coach)
	}
	return coaches, rows.Err()
}

func (r *coachRepository) Update(ctx context.Context, coach *model.Coach) error {
	query := `UPDATE coaches SET name = ?, timezone = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, coach.Name, coach.Timezone, coach.ID)
	return err
}

func (r *coachRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE coaches SET deleted = TRUE WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
