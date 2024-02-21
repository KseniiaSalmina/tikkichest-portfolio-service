package connector

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/models"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/storage/postgresql"
)

type PostgresConnector struct {
	db *pgxpool.Pool
}

func NewPostgresConnector(db *pgxpool.Pool) *PostgresConnector {
	return &PostgresConnector{db: db}
}

func (pc *PostgresConnector) GetAllPortfolios(ctx context.Context, limit int, offset int, filter postgresql.PortfoliosFilter) ([]models.Portfolio, int, error) {
	portfolios, err := postgresql.GetAllPortfolios(ctx, pc.db, limit, offset, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get data from database: %w", err)
	}

	rowsAmount, err := postgresql.CountPortfoliosPages(ctx, pc.db, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pages: %w", err)
	}

	var pageAmount int
	if rowsAmount%limit != 0 {
		pageAmount = rowsAmount/limit + 1
	} else {
		pageAmount = rowsAmount / limit
	}

	return portfolios, pageAmount, nil
}

func (pc *PostgresConnector) GetPortfolioByID(ctx context.Context, portfolioID int) (*models.Portfolio, error) {
	return postgresql.GetPortfolioByID(ctx, pc.db, portfolioID)
}

func (pc *PostgresConnector) CreatePortfolio(ctx context.Context, portfolio models.Portfolio) (int, error) {
	return postgresql.CreatePortfolio(ctx, pc.db, portfolio)
}

func (pc *PostgresConnector) PatchPortfolio(ctx context.Context, portfolio models.Portfolio) error {
	return postgresql.PatchPortfolio(ctx, pc.db, portfolio)
}

func (pc *PostgresConnector) DeletePortfolio(ctx context.Context, portfolioID int) error {
	return postgresql.DeletePortfolio(ctx, pc.db, portfolioID)
}

func (pc *PostgresConnector) CreateCategory(ctx context.Context, name string) (int, error) {
	return postgresql.CreateCategory(ctx, pc.db, name)
}

func (pc *PostgresConnector) DeleteCategory(ctx context.Context, id int) error {
	return postgresql.DeleteCategory(ctx, pc.db, id)
}

func (pc *PostgresConnector) GetAllCategories(ctx context.Context, limit int, offset int) ([]models.Category, int, error) {
	categories, err := postgresql.GetAllCategories(ctx, pc.db, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get data from database: %w", err)
	}

	rowsAmount, err := postgresql.CountCategoriesPages(ctx, pc.db)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pages: %w", err)
	}

	var pageAmount int
	if rowsAmount%limit != 0 {
		pageAmount = rowsAmount/limit + 1
	} else {
		pageAmount = rowsAmount / limit
	}

	return categories, pageAmount, nil
}

func (pc *PostgresConnector) GetAllCraftsByPortfolioID(ctx context.Context, portfolioID int, limit int, offset int) ([]models.Craft, int, error) {
	crafts, err := postgresql.GetAllCraftsByPortfolioID(ctx, pc.db, portfolioID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get data from database: %w", err)
	}

	rowsAmount, err := postgresql.CountCraftsPages(ctx, pc.db, portfolioID, true)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pages: %w", err)
	}

	var pageAmount int
	if rowsAmount%limit != 0 {
		pageAmount = rowsAmount/limit + 1
	} else {
		pageAmount = rowsAmount / limit
	}

	return crafts, pageAmount, nil
}

func (pc *PostgresConnector) GetCraftByID(ctx context.Context, craftID int) (*models.Craft, error) {
	return postgresql.GetCraftByID(ctx, pc.db, craftID)
}

func (pc *PostgresConnector) CreateCraft(ctx context.Context, portfolioID int, craft models.Craft) (int, error) {
	return postgresql.CreateCraft(ctx, pc.db, portfolioID, craft)
}

func (pc *PostgresConnector) AddTagToCraft(ctx context.Context, craftID int, tagID int) error {
	return postgresql.AddTagToCraft(ctx, pc.db, craftID, tagID)
}

func (pc *PostgresConnector) DeleteTagFromCraft(ctx context.Context, craftID int, tagID int) error {
	return postgresql.DeleteTagFromCraft(ctx, pc.db, craftID, tagID)
}

func (pc *PostgresConnector) PatchCraft(ctx context.Context, craft models.Craft) error {
	return postgresql.PatchCraft(ctx, pc.db, craft)
}

func (pc *PostgresConnector) DeleteCraft(ctx context.Context, id int) error {
	return postgresql.DeleteCraft(ctx, pc.db, id)
}

func (pc *PostgresConnector) GetAllCraftsByTagID(ctx context.Context, tagID int, limit int, offset int) ([]models.Craft, int, error) {
	crafts, err := postgresql.GetAllCraftsByTagID(ctx, pc.db, tagID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get data from database: %w", err)
	}

	rowsAmount, err := postgresql.CountCraftsPages(ctx, pc.db, tagID, false)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pages: %w", err)
	}

	var pageAmount int
	if rowsAmount%limit != 0 {
		pageAmount = rowsAmount/limit + 1
	} else {
		pageAmount = rowsAmount / limit
	}

	return crafts, pageAmount, nil
}

func (pc *PostgresConnector) GetAllTags(ctx context.Context, limit int, offset int) ([]models.Tag, int, error) {
	tags, err := postgresql.GetAllTags(ctx, pc.db, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get data from database: %w", err)
	}

	rowsAmount, err := postgresql.CountTagsPages(ctx, pc.db)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pages: %w", err)
	}

	var pageAmount int
	if rowsAmount%limit != 0 {
		pageAmount = rowsAmount/limit + 1
	} else {
		pageAmount = rowsAmount / limit
	}

	return tags, pageAmount, nil
}

func (pc *PostgresConnector) CreateTag(ctx context.Context, name string) (int, error) {
	return postgresql.CreateTag(ctx, pc.db, name)
}

func (pc *PostgresConnector) DeleteTag(ctx context.Context, id int) error {
	return postgresql.DeleteTag(ctx, pc.db, id)
}

func (pc *PostgresConnector) CreateContent(ctx context.Context, craftID int, content models.Content) (int, error) {
	return postgresql.CreateContent(ctx, pc.db, craftID, content)
}

func (pc *PostgresConnector) DeleteContent(ctx context.Context, id int) error {
	return postgresql.DeleteContent(ctx, pc.db, id)
}

func (pc *PostgresConnector) PatchContent(ctx context.Context, content models.Content) error {
	return postgresql.PatchContent(ctx, pc.db, content)
}
