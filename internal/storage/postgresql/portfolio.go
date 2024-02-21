package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strconv"

	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/models"
)

func CreatePortfolio(ctx context.Context, db *pgxpool.Pool, portfolio models.Portfolio) (int, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to create portfolio: transaction error: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Println(err)
		}
	}() // TODO: прикрутить логгер или убрать логирование ошибки

	var id pgtype.Int8
	if err = tx.QueryRow(ctx, `INSERT INTO portfolios (profile_id, name, category_id, description) VALUES ($1, $2, $3, $4) RETURNING id`, portfolio.ProfileID, portfolio.Name, portfolio.Category.ID, portfolio.Description).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to create portfolio: creation error: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to create portfolio: transaction error: %w", err)
	}

	return int(id.Int), nil
}

func CreateCategory(ctx context.Context, db *pgxpool.Pool, name string) (int, error) {
	var id pgtype.Int8
	if err := db.QueryRow(ctx, `INSERT INTO categories (name) VALUES ($1) RETURNING id`, name).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to create category: %w", err)
	}

	return int(id.Int), nil
}

func DeleteCategory(ctx context.Context, db *pgxpool.Pool, id int) error {
	if _, err := db.Exec(ctx, `DELETE FROM categories WHERE id=$1`, id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func GetAllCategories(ctx context.Context, db *pgxpool.Pool, limit, offset int) ([]models.Category, error) {
	rows, err := db.Query(ctx, `SELECT id, name FROM categories LIMIT $1 OFFSET $2`, limit, offset)
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

func CountCategoriesPages(ctx context.Context, db *pgxpool.Pool) (int, error) {
	var amount pgtype.Int8
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM categories`).Scan(&amount); err != nil {
		return 0, fmt.Errorf("failed to count categories: %w", err)
	}

	return int(amount.Int), nil
}

func DeletePortfolio(ctx context.Context, db *pgxpool.Pool, portfolioID int) error {
	if _, err := db.Exec(ctx, `DELETE FROM portfolios WHERE id=$1`, portfolioID); err != nil {
		return fmt.Errorf("failed to delete portfolio: %w", err)
	}

	return nil
}

func PatchPortfolio(ctx context.Context, db *pgxpool.Pool, portfolio models.Portfolio) error {
	if _, err := db.Exec(ctx, `UPDATE portfolios SET name = $1, description = $2, category_id = $3 WHERE id = $4`, portfolio.Name, portfolio.Description, portfolio.Category.ID, portfolio.ID); err != nil {
		return fmt.Errorf("failed to update portfolio: %w", err)
	}

	return nil
}

func GetPortfolioByID(ctx context.Context, db *pgxpool.Pool, portfolioID int) (*models.Portfolio, error) {
	var profileID, categoryID pgtype.Int8
	var portfolioName, categoryName, portfolioDescription pgtype.Text

	if err := db.QueryRow(ctx, `SELECT portfolios.profile_id, portfolios.name, portfolios.category_id, categories.name, portfolios.description 
FROM portfolios JOIN categories ON portfolios.category_id = categories.id 
WHERE portfolios.id = $1`, portfolioID).Scan(&profileID, &portfolioName, &categoryID, &categoryName, &portfolioDescription); err != nil {
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}

	return &models.Portfolio{ID: portfolioID, ProfileID: int(profileID.Int), Name: portfolioName.String, Description: portfolioDescription.String, Category: models.Category{ID: int(categoryID.Int), Name: categoryName.String}}, nil
}

type PortfoliosFilter struct {
	Type PortfoliosFilterType
	ID   int
}

type PortfoliosFilterType string

const (
	ByProfileID  PortfoliosFilterType = "portfolios.profile_id ="
	ByCategoryID PortfoliosFilterType = "portfolios.category_id ="
)

func GetAllPortfolios(ctx context.Context, db *pgxpool.Pool, limit, offset int, filter PortfoliosFilter) ([]models.Portfolio, error) {
	var sql string
	switch {
	case filter.Type == "":
		sql = `SELECT portfolios.id, portfolios.profile_id, portfolios.name, portfolios.description, portfolios.category_id, categories.name 
FROM portfolios JOIN categories ON portfolios.category_id = categories.id`
	case filter.Type == ByProfileID || filter.Type == ByCategoryID:
		sql = `SELECT portfolios.id, portfolios.profile_id, portfolios.name, portfolios.description, portfolios.category_id, categories.name 
FROM portfolios JOIN categories ON portfolios.category_id = categories.id WHERE ` + string(filter.Type) + strconv.Itoa(filter.ID)
	default:
		return nil, errors.New("failed to get portfolios: incorrect filter")
	}

	rows, err := db.Query(ctx, sql+` LIMIT $1 OFFSET $2`, limit, offset)
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

func CountPortfoliosPages(ctx context.Context, db *pgxpool.Pool, filter PortfoliosFilter) (int, error) {
	var sql string
	switch {
	case filter.Type == "":
		sql = `SELECT COUNT(*) FROM categories`
	case filter.Type == ByProfileID || filter.Type == ByCategoryID:
		sql = `SELECT COUNT(*) FROM categories WHERE ` + string(filter.Type) + strconv.Itoa(filter.ID)
	default:
		return 0, errors.New("failed to count portfolios: incorrect filter")
	}

	var amount pgtype.Int8
	if err := db.QueryRow(ctx, sql).Scan(&amount); err != nil {
		return 0, fmt.Errorf("failed to count portfolios: %w", err)
	}

	return int(amount.Int), nil
}
