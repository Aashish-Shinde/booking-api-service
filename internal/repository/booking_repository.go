package repository

import (
	"context"
	"database/sql"

	"github.com/booking-api-service/internal/model"
)

type bookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, booking *model.Booking) (int64, error) {
	query := `INSERT INTO bookings (user_id, coach_id, start_time, end_time, status, idempotency_key) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, booking.UserID, booking.CoachID, booking.StartTime, booking.EndTime, booking.Status, booking.IdempotencyKey)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *bookingRepository) GetByID(ctx context.Context, id int64) (*model.Booking, error) {
	query := `SELECT id, user_id, coach_id, start_time, end_time, status, idempotency_key, created_at, updated_at, deleted 
	          FROM bookings WHERE id = ? AND deleted = FALSE`
	booking := &model.Booking{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&booking.ID, &booking.UserID, &booking.CoachID, &booking.StartTime, &booking.EndTime,
		&booking.Status, &booking.IdempotencyKey, &booking.CreatedAt, &booking.UpdatedAt, &booking.Deleted,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return booking, nil
}

func (r *bookingRepository) GetByUser(ctx context.Context, userID int64) ([]*model.Booking, error) {
	query := `SELECT id, user_id, coach_id, start_time, end_time, status, idempotency_key, created_at, updated_at, deleted 
	          FROM bookings WHERE user_id = ? AND deleted = FALSE ORDER BY start_time DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*model.Booking
	for rows.Next() {
		b := &model.Booking{}
		err := rows.Scan(&b.ID, &b.UserID, &b.CoachID, &b.StartTime, &b.EndTime, &b.Status, &b.IdempotencyKey, &b.CreatedAt, &b.UpdatedAt, &b.Deleted)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, rows.Err()
}

func (r *bookingRepository) GetByCoachAndDateRange(ctx context.Context, coachID int64, startTime string, endTime string) ([]*model.Booking, error) {
	query := `SELECT id, user_id, coach_id, start_time, end_time, status, idempotency_key, created_at, updated_at, deleted 
	          FROM bookings WHERE coach_id = ? AND start_time >= ? AND end_time <= ? AND status = ? AND deleted = FALSE`
	rows, err := r.db.QueryContext(ctx, query, coachID, startTime, endTime, model.BookingStatusActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*model.Booking
	for rows.Next() {
		b := &model.Booking{}
		err := rows.Scan(&b.ID, &b.UserID, &b.CoachID, &b.StartTime, &b.EndTime, &b.Status, &b.IdempotencyKey, &b.CreatedAt, &b.UpdatedAt, &b.Deleted)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, rows.Err()
}

func (r *bookingRepository) Update(ctx context.Context, booking *model.Booking) error {
	query := `UPDATE bookings SET user_id = ?, coach_id = ?, start_time = ?, end_time = ?, status = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, booking.UserID, booking.CoachID, booking.StartTime, booking.EndTime, booking.Status, booking.ID)
	return err
}

func (r *bookingRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE bookings SET deleted = TRUE WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *bookingRepository) GetByIdempotencyKey(ctx context.Context, key string) (*model.Booking, error) {
	query := `SELECT id, user_id, coach_id, start_time, end_time, status, idempotency_key, created_at, updated_at, deleted 
	          FROM bookings WHERE idempotency_key = ? AND deleted = FALSE`
	booking := &model.Booking{}
	err := r.db.QueryRowContext(ctx, query, key).Scan(
		&booking.ID, &booking.UserID, &booking.CoachID, &booking.StartTime, &booking.EndTime,
		&booking.Status, &booking.IdempotencyKey, &booking.CreatedAt, &booking.UpdatedAt, &booking.Deleted,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return booking, nil
}
