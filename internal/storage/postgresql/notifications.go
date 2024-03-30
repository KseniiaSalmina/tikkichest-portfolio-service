package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5"
)

func (db *DB) NotificationsOn(ctx context.Context, userID int) error {
	if _, err := db.db.Exec(ctx, `INSERT INTO notifications (profile_id) VALUES ($1)`, userID); err != nil {
		return fmt.Errorf("failed to save userID: %w", err)
	}

	return nil
}

func (db *DB) NotificationsOff(ctx context.Context, userID int) error {
	if _, err := db.db.Exec(ctx, `DELETE FROM notifications WHERE profile_id = $1`, userID); err != nil {
		return fmt.Errorf("failed to delete userID: %w", err)
	}

	return nil
}

func (db *DB) IsNotificationsOn(ctx context.Context, userID int) (bool, error) {
	var id pgtype.Int4
	if err := db.db.QueryRow(ctx, `SELECT profile_id FROM notifications WHERE profile_id = $1`, userID).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get userID: %w", err)
	}

	return true, nil
}
