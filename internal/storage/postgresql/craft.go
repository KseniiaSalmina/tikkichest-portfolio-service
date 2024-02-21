package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"

	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/models"
)

func CreateCraft(ctx context.Context, db *pgxpool.Pool, portfolioID int, craft models.Craft) (int, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to create craft: transaction error: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Println(err)
		}
	}() //TODO: прикрутить логгер или убрать логирование ошибки

	var craftID pgtype.Int8
	if err = tx.QueryRow(ctx, `INSERT INTO crafts (portfolio_id, name, description) VALUES ($1, $2, $3) RETURNING id`, portfolioID, craft.Name, craft.Description).Scan(&craftID); err != nil {
		return 0, fmt.Errorf("failed to create craft: %w", err)
	}

	for _, tag := range craft.Tags {
		if _, err = tx.Exec(ctx, `INSERT INTO crafts_tags (craft_id, tag_id) VALUES ($1, $2)`, int(craftID.Int), tag.ID); err != nil {
			return 0, fmt.Errorf("failed to create craft: tags error: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to create craft: transaction error: %w", err)
	}

	return int(craftID.Int), nil
}

func CreateTag(ctx context.Context, db *pgxpool.Pool, name string) (int, error) {
	var id pgtype.Int8
	if err := db.QueryRow(ctx, `INSERT INTO tags (name) VALUES ($1) RETURNING id`, name).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to create tag: %w", err)
	}

	return int(id.Int), nil
}

func DeleteTag(ctx context.Context, db *pgxpool.Pool, id int) error {
	if _, err := db.Exec(ctx, `DELETE FROM tags WHERE id = $1`, id); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

func GetAllTags(ctx context.Context, db *pgxpool.Pool, limit, offset int) ([]models.Tag, error) {
	rows, err := db.Query(ctx, `SELECT id, name FROM tags LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all tags: %w", err)
	}

	var tags []models.Tag

	for rows.Next() {
		var id pgtype.Int8
		var name pgtype.Text

		if err = rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("failed to get all tags: scan error: %w", err)
		}

		tag := models.Tag{ID: int(id.Int), Name: name.String}
		tags = append(tags, tag)
	}

	return tags, nil
}

func CountTagsPages(ctx context.Context, db *pgxpool.Pool) (int, error) {
	var amount pgtype.Int8

	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM tags`).Scan(&amount); err != nil {
		return 0, fmt.Errorf("failed to count tags: %w", err)
	}

	return int(amount.Int), nil
}

func AddTagToCraft(ctx context.Context, db *pgxpool.Pool, craftID, tagID int) error {
	if _, err := db.Exec(ctx, `INSERT INTO crafts_tags (craft_id, tag_id) VALUES ($1, $2)`, craftID, tagID); err != nil {
		return fmt.Errorf("failed to add tag: %w", err)
	}

	return nil
}

func DeleteTagFromCraft(ctx context.Context, db *pgxpool.Pool, craftID, tagID int) error {
	if _, err := db.Exec(ctx, `DELETE FROM crafts_tags WHERE craft_id=$1 AND tag_id=$2`, craftID, tagID); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

func DeleteCraft(ctx context.Context, db *pgxpool.Pool, id int) error {
	if _, err := db.Exec(ctx, `DELETE FROM crafts WHERE id=$1`, id); err != nil {
		return fmt.Errorf("failed to delete craft: %w", err)
	}

	return nil
}

func PatchCraft(ctx context.Context, db *pgxpool.Pool, craft models.Craft) error {
	if _, err := db.Exec(ctx, `UPDATE crafts SET name = $1, description = $2 WHERE id = $3`, craft.Name, craft.Description, craft.ID); err != nil {
		return fmt.Errorf("failed to update craft: %w", err)
	}

	return nil
}

func GetCraftByID(ctx context.Context, db *pgxpool.Pool, craftID int) (*models.Craft, error) {
	craft := models.Craft{ID: craftID}

	var craftName, craftDescription pgtype.Text
	if err := db.QueryRow(ctx, `SELECT name, description FROM crafts WHERE id = $1`, craftID).Scan(&craftName, &craftDescription); err != nil {
		return nil, fmt.Errorf("failed to get craft: %w", err)
	}
	craft.Name, craft.Description = craftName.String, craftDescription.String

	rows, err := db.Query(ctx, `SELECT crafts_tags.tag_id, tags.name FROM crafts_tags JOIN tags ON crafts_tags.tag_id = tags.id WHERE crafts_tags.craft_id = $1`, craftID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to get craft: tags error: %w", err)
	}

	var tag models.Tag
	for rows.Next() {
		var tagID pgtype.Int8
		var tagName pgtype.Text

		if err = rows.Scan(&tagID, &tagName); err != nil {
			return nil, fmt.Errorf("failed to get craft: scan tags error: %w", err)
		}

		tag.ID, tag.Name = int(tagID.Int), tagName.String
		craft.Tags = append(craft.Tags, tag)
	}

	rows, err = db.Query(ctx, `SELECT id, description, data FROM contents WHERE craft_id = $1`, craftID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to get craft: content error: %w", err)
		}
		return &craft, nil
	}

	for rows.Next() {
		var contentID pgtype.Int8
		var contentDescription pgtype.Text
		var contentData pgtype.Bytea

		if err = rows.Scan(&contentID, &contentDescription, &contentData); err != nil {
			return nil, fmt.Errorf("failed to get craft: content scan error: %w", err)
		}

		content := models.Content{ID: int(contentID.Int), Description: contentDescription.String, Data: contentData.Bytes}
		craft.Contents = append(craft.Contents, content)
	}

	return &craft, nil
}

func GetAllCraftsByPortfolioID(ctx context.Context, db *pgxpool.Pool, portfolioID, limit, offset int) ([]models.Craft, error) {
	var crafts []models.Craft

	rows, err := db.Query(ctx, `SELECT id, name, description FROM crafts WHERE portfolio_id = $1 LIMIT $2 OFFSET $3`, portfolioID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get crafts by portfolio id: %w", err)
	}

	for rows.Next() {
		var craftID pgtype.Int8
		var craftName, craftDescription pgtype.Text

		if err = rows.Scan(&craftID, &craftName, &craftDescription); err != nil {
			return nil, fmt.Errorf("failed to get crafts: scan error %w", err)
		}

		craft := models.Craft{ID: int(craftID.Int), Name: craftName.String, Description: craftDescription.String}
		crafts = append(crafts, craft)
	}

	for i, craft := range crafts {
		tags, content, err := getDetailsOfCraft(ctx, db, craft.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get crafts: details error: %w", err)
		}

		crafts[i].Tags = tags
		crafts[i].Contents = []models.Content{content}
	}

	return crafts, nil
}

func CountCraftsPages(ctx context.Context, db *pgxpool.Pool, id int, isPortfolioID bool) (int, error) {
	var amount pgtype.Int8
	var err error

	if isPortfolioID {
		err = db.QueryRow(ctx, `SELECT COUNT(*) FROM crafts WHERE portfolio_id = $1`, id).Scan(&amount)
	} else {
		err = db.QueryRow(ctx, `SELECT COUNT(*) FROM crafts_tags WHERE tag_id = $1`, id).Scan(&amount)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to count crafts: %w", err)
	}

	return int(amount.Int), nil
}

func getDetailsOfCraft(ctx context.Context, db *pgxpool.Pool, craftID int) ([]models.Tag, models.Content, error) {
	rows, err := db.Query(ctx, `SELECT crafts_tags.tag_id, tags.name FROM crafts_tags JOIN tags ON crafts_tags.tag_id = tags.id WHERE crafts_tags.craft_id = $1`, craftID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, models.Content{}, fmt.Errorf("failed to get details of craft: tags error: %w", err)
	}

	var tags []models.Tag
	for rows.Next() {
		var tagID pgtype.Int8
		var tagName pgtype.Text

		if err = rows.Scan(&tagID, &tagName); err != nil {
			return nil, models.Content{}, fmt.Errorf("failed to get details of craft: scan tags error: %w", err)
		}

		tag := models.Tag{ID: int(tagID.Int), Name: tagName.String}
		tags = append(tags, tag)
	}

	var contentID pgtype.Int8
	var contentDescription pgtype.Text
	var contentData pgtype.Bytea

	if err = db.QueryRow(ctx, `SELECT id, description, data FROM contents WHERE craft_id = $1 LIMIT 1`, craftID).Scan(&contentID, &contentDescription, &contentData); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, models.Content{}, fmt.Errorf("failed to get details of craft: %w", err)
		}
		return tags, models.Content{}, nil
	}

	content := models.Content{ID: int(contentID.Int), Description: contentDescription.String, Data: contentData.Bytes}

	return tags, content, nil
}

func GetAllCraftsByTagID(ctx context.Context, db *pgxpool.Pool, tagID, limit, offset int) ([]models.Craft, error) {
	rows, err := db.Query(ctx, `SELECT craft_id FROM crafts_tags WHERE tag_id = $1 LIMIT $2 OFFSET $3`, tagID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get crafts by tag id: %w", err)
	}

	var craftsIDs []int
	for rows.Next() {
		var craftID pgtype.Int8

		if err = rows.Scan(&craftID); err != nil {
			return nil, fmt.Errorf("failed to get crafts by portfolio id: scan craftID error %w", err)
		}

		craftsIDs = append(craftsIDs, int(craftID.Int))
	}

	var crafts []models.Craft

	for _, craftID := range craftsIDs {
		var craftName, craftDescription pgtype.Text
		if err = db.QueryRow(ctx, `SELECT name, description FROM crafts WHERE id = $1`, craftID).Scan(&craftName, &craftDescription); err != nil {
			return nil, fmt.Errorf("failed to get crafts by portfolio id: scan craft error %w", err)
		}

		tags, content, err := getDetailsOfCraft(ctx, db, craftID)
		if err != nil {
			return nil, fmt.Errorf("failed to get crafts: details error: %w", err)
		}

		craft := models.Craft{ID: craftID, Name: craftName.String, Description: craftDescription.String, Tags: tags, Contents: []models.Content{content}}
		crafts = append(crafts, craft)
	}

	return crafts, nil
}
