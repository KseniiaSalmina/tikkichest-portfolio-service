package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/pgtype"

	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/models"
)

func (db *DB) CreateContent(ctx context.Context, craftID int, content models.Content) (int, error) {
	var contentID pgtype.Int8
	if err := db.db.QueryRow(ctx, `INSERT INTO contents (craft_id, description, data) VALUES ($1, $2, $3) RETURNING id`, craftID, content.Description, content.Data).Scan(&contentID); err != nil {
		return 0, fmt.Errorf("failed to create content: %w", err)
	}

	return int(contentID.Int), nil
}

func (db *DB) DeleteContent(ctx context.Context, id int) error {
	if _, err := db.db.Exec(ctx, `DELETE FROM contents WHERE id = $1`, id); err != nil {
		return fmt.Errorf("failed to delete content: %w", err)
	}

	return nil
}

func (db *DB) PatchContent(ctx context.Context, content models.Content) error {
	if _, err := db.db.Exec(ctx, `UPDATE contents SET description = $1, data = $2 WHERE id = $3`, content.Description, content.Data, content.ID); err != nil {
		return fmt.Errorf("failed to update content: %w", err)
	}

	return nil
}
