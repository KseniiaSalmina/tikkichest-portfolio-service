package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/pgtype"
	"strings"

	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/models"
)

func (db *DB) CreatePortfolio(ctx context.Context, portfolio models.Portfolio) (int, error) {
	tx, err := db.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to create portfolio: transaction error: %w", err)
	}

	defer tx.Rollback(ctx)

	var id pgtype.Int8
	if err = tx.QueryRow(ctx, `INSERT INTO portfolios (profile_id, name, category_id, description) VALUES ($1, $2, $3, $4) RETURNING id`, portfolio.ProfileID, portfolio.Name, portfolio.Category.ID, portfolio.Description).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to create portfolio: creation error: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to create portfolio: transaction error: %w", err)
	}

	return int(id.Int), nil
}

func (db *DB) CreateCategory(ctx context.Context, name string) (int, error) {
	var id pgtype.Int8
	if err := db.db.QueryRow(ctx, `INSERT INTO categories (name) VALUES ($1) RETURNING id`, name).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to create category: %w", err)
	}

	return int(id.Int), nil
}

func (db *DB) DeleteCategory(ctx context.Context, id int) error {
	if _, err := db.db.Exec(ctx, `DELETE FROM categories WHERE id=$1`, id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (db *DB) GetAllCategories(ctx context.Context, limit, offset int) ([]models.Category, error) {
	rows, err := db.db.Query(ctx, `SELECT id, name FROM categories LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all categories: %w", err)
	}

	var categories []models.Category
	var category models.Category

	for rows.Next() {
		var id pgtype.Int8
		var name pgtype.Text
		if err = rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("failed to get all categories: scan error: %w", err)
		}
		category.ID, category.Name = int(id.Int), name.String
		categories = append(categories, category)
	}

	return categories, nil
}

func (db *DB) CountCategoriesPages(ctx context.Context) (int, error) {
	var amount pgtype.Int8
	if err := db.db.QueryRow(ctx, `SELECT COUNT(*) FROM categories`).Scan(&amount); err != nil {
		return 0, fmt.Errorf("failed to count categories: %w", err)
	}

	return int(amount.Int), nil
}

func (db *DB) DeletePortfolio(ctx context.Context, portfolioID int) error {
	if _, err := db.db.Exec(ctx, `DELETE FROM portfolios WHERE id=$1`, portfolioID); err != nil {
		return fmt.Errorf("failed to delete portfolio: %w", err)
	}

	return nil
}

func (db *DB) PatchPortfolio(ctx context.Context, portfolio models.Portfolio) error {
	if _, err := db.db.Exec(ctx, `UPDATE portfolios SET name = $1, description = $2, category_id = $3 WHERE id = $4`, portfolio.Name, portfolio.Description, portfolio.Category.ID, portfolio.ID); err != nil {
		return fmt.Errorf("failed to update portfolio: %w", err)
	}

	return nil
}

func (db *DB) GetPortfolioByID(ctx context.Context, portfolioID int) (*models.Portfolio, error) {
	var profileID, categoryID pgtype.Int8
	var portfolioName, categoryName, portfolioDescription pgtype.Text

	if err := db.db.QueryRow(ctx, `
	SELECT portfolios.profile_id, 
       portfolios.name, 
       portfolios.category_id, 
       categories.name, 
       portfolios.description 
	FROM portfolios 
	JOIN categories ON portfolios.category_id = categories.id 
	WHERE portfolios.id = $1`,
		portfolioID).Scan(&profileID, &portfolioName, &categoryID, &categoryName, &portfolioDescription); err != nil {
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}

	return &models.Portfolio{ID: portfolioID, ProfileID: int(profileID.Int), Name: portfolioName.String, Description: portfolioDescription.String, Category: models.Category{ID: int(categoryID.Int), Name: categoryName.String}}, nil
}

func (db *DB) GetAllPortfolios(ctx context.Context, limit, offset int, id int, filterType PortfoliosFilterType) ([]models.Portfolio, error) {
	filter, err := portfolioFilter(filterType, id)
	if err != nil {
		return nil, err
	}

	sql := strings.Join([]string{`
	SELECT portfolios.id, 
       portfolios.profile_id, 
       portfolios.name, 
       portfolios.description, 
       portfolios.category_id, 
       categories.name 
	FROM portfolios 
    JOIN categories ON portfolios.category_id = categories.id`,
		filter,
		`LIMIT $1 OFFSET $2`}, " ")

	rows, err := db.db.Query(ctx, sql, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolios: %w", err)
	}

	var portfolios []models.Portfolio
	for rows.Next() {
		var portfolioID, profileID, categoryID pgtype.Int8
		var portfolioName, categoryName, portfolioDescription pgtype.Text

		if err = rows.Scan(&portfolioID, &profileID, &portfolioName, &portfolioDescription, &categoryID, &categoryName); err != nil {
			return nil, fmt.Errorf("failed to get portfolios: scan error: %w", err)
		}

		portfolio := models.Portfolio{ID: int(portfolioID.Int), ProfileID: int(profileID.Int), Name: portfolioName.String, Description: portfolioDescription.String, Category: models.Category{ID: int(categoryID.Int), Name: categoryName.String}}
		portfolios = append(portfolios, portfolio)
	}

	return portfolios, nil
}

func (db *DB) CountPortfoliosPages(ctx context.Context, id int, filterType PortfoliosFilterType) (int, error) {
	filter, err := portfolioFilter(filterType, id)
	if err != nil {
		return 0, err
	}

	sql := strings.Join([]string{"SELECT COUNT(*) FROM categories", filter}, " ")

	var amount pgtype.Int8
	if err := db.db.QueryRow(ctx, sql).Scan(&amount); err != nil {
		return 0, fmt.Errorf("failed to count portfolios: %w", err)
	}

	return int(amount.Int), nil
}
