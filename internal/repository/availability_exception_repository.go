package repository

import (
	"context"
	"database/sql"

	"github.com/booking-api-service/internal/model"
)

type availabilityExceptionRepository struct {
	db *sql.DB
}

func NewAvailabilityExceptionRepository(db *sql.DB) AvailabilityExceptionRepository {
	return &availabilityExceptionRepository{db: db}
}

func (r *availabilityExceptionRepository) Create(ctx context.Context, exception *model.AvailabilityException) (int64, error) {
	query := `INSERT INTO availability_exceptions (coach_id, date, start_time, end_time, is_available) 
	          VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, exception.CoachID, exception.Date, exception.StartTime, exception.EndTime, exception.IsAvailable)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *availabilityExceptionRepository) GetByCoachAndDate(ctx context.Context, coachID int64, date string) (*model.AvailabilityException, error) {
	query := `SELECT id, coach_id, date, start_time, end_time, is_available, created_at, updated_at, deleted 
	          FROM availability_exceptions WHERE coach_id = ? AND DATE(date) = ? AND deleted = FALSE`
	exception := &model.AvailabilityException{}
	err := r.db.QueryRowContext(ctx, query, coachID, date).Scan(
		&exception.ID, &exception.CoachID, &exception.Date, &exception.StartTime, &exception.EndTime,
		&exception.IsAvailable, &exception.CreatedAt, &exception.UpdatedAt, &exception.Deleted,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return exception, nil
}

func (r *availabilityExceptionRepository) GetByCoach(ctx context.Context, coachID int64) ([]*model.AvailabilityException, error) {
	query := `SELECT id, coach_id, date, start_time, end_time, is_available, created_at, updated_at, deleted 
	          FROM availability_exceptions WHERE coach_id = ? AND deleted = FALSE`
	rows, err := r.db.QueryContext(ctx, query, coachID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exceptions []*model.AvailabilityException
	for rows.Next() {
		e := &model.AvailabilityException{}
		err := rows.Scan(&e.ID, &e.CoachID, &e.Date, &e.StartTime, &e.EndTime, &e.IsAvailable, &e.CreatedAt, &e.UpdatedAt, &e.Deleted)
		if err != nil {
			return nil, err
		}
		exceptions = append(exceptions, e)
	}
	return exceptions, rows.Err()
}

func (r *availabilityExceptionRepository) Update(ctx context.Context, exception *model.AvailabilityException) error {
	query := `UPDATE availability_exceptions SET start_time = ?, end_time = ?, is_available = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, exception.StartTime, exception.EndTime, exception.IsAvailable, exception.ID)
	return err
}

func (r *availabilityExceptionRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE availability_exceptions SET deleted = TRUE WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
