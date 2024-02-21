package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/models"
)

func CreateContent(ctx context.Context, db *pgxpool.Pool, craftID int, content models.Content) (int, error) {
	var contentID pgtype.Int8
	if err := db.QueryRow(ctx, `INSERT INTO contents (craft_id, description, data) VALUES ($1, $2, $3) RETURNING id`, craftID, content.Description, content.Data).Scan(&contentID); err != nil {
		return 0, fmt.Errorf("failed to create content: %w", err)
	}

	return int(contentID.Int), nil
}

func DeleteContent(ctx context.Context, db *pgxpool.Pool, id int) error {
	if _, err := db.Exec(ctx, `DELETE FROM contents WHERE id = $1`, id); err != nil {
		return fmt.Errorf("failed to delete content: %w", err)
	}

	return nil
}

func PatchContent(ctx context.Context, db *pgxpool.Pool, content models.Content) error {
	if _, err := db.Exec(ctx, `UPDATE contents SET description = $1, data = $2 WHERE id = $3`, content.Description, content.Data, content.ID); err != nil {
		return fmt.Errorf("failed to update content: %w", err)
	}

	return nil
}
