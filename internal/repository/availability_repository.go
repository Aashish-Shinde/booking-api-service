package repository

import (
	"context"
	"database/sql"

	"github.com/booking-api-service/internal/model"
)

type availabilityRepository struct {
	db *sql.DB
}

func NewAvailabilityRepository(db *sql.DB) AvailabilityRepository {
	return &availabilityRepository{db: db}
}

func (r *availabilityRepository) Create(ctx context.Context, availability *model.Availability) (int64, error) {
	query := `INSERT INTO availability (coach_id, day_of_week, start_time, end_time) VALUES (?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, availability.CoachID, availability.DayOfWeek, availability.StartTime, availability.EndTime)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *availabilityRepository) GetByCoachAndDay(ctx context.Context, coachID int64, dayOfWeek int) ([]*model.Availability, error) {
	query := `SELECT id, coach_id, day_of_week, start_time, end_time, created_at, updated_at, deleted 
	          FROM availability WHERE coach_id = ? AND day_of_week = ? AND deleted = FALSE`
	rows, err := r.db.QueryContext(ctx, query, coachID, dayOfWeek)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var availabilities []*model.Availability
	for rows.Next() {
		a := &model.Availability{}
		err := rows.Scan(&a.ID, &a.CoachID, &a.DayOfWeek, &a.StartTime, &a.EndTime, &a.CreatedAt, &a.UpdatedAt, &a.Deleted)
		if err != nil {
			return nil, err
		}
		availabilities = append(availabilities, a)
	}
	return availabilities, rows.Err()
}

func (r *availabilityRepository) GetByCoach(ctx context.Context, coachID int64) ([]*model.Availability, error) {
	query := `SELECT id, coach_id, day_of_week, start_time, end_time, created_at, updated_at, deleted 
	          FROM availability WHERE coach_id = ? AND deleted = FALSE`
	rows, err := r.db.QueryContext(ctx, query, coachID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var availabilities []*model.Availability
	for rows.Next() {
		a := &model.Availability{}
		err := rows.Scan(&a.ID, &a.CoachID, &a.DayOfWeek, &a.StartTime, &a.EndTime, &a.CreatedAt, &a.UpdatedAt, &a.Deleted)
		if err != nil {
			return nil, err
		}
		availabilities = append(availabilities, a)
	}
	return availabilities, rows.Err()
}

func (r *availabilityRepository) Update(ctx context.Context, availability *model.Availability) error {
	query := `UPDATE availability SET coach_id = ?, day_of_week = ?, start_time = ?, end_time = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, availability.CoachID, availability.DayOfWeek, availability.StartTime, availability.EndTime, availability.ID)
	return err
}

func (r *availabilityRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE availability SET deleted = TRUE WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *availabilityRepository) DeleteByCoachAndDay(ctx context.Context, coachID int64, dayOfWeek int) error {
	query := `UPDATE availability SET deleted = TRUE WHERE coach_id = ? AND day_of_week = ?`
	_, err := r.db.ExecContext(ctx, query, coachID, dayOfWeek)
	return err
}
