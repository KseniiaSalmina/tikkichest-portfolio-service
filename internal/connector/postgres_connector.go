package connector

import (
	"context"
	"fmt"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/models"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/storage/postgresql"
)

type PostgresConnector struct {
	db *postgresql.DB
}

func NewPostgresConnector(db *postgresql.DB) *PostgresConnector {
	return &PostgresConnector{db: db}
}

func (pc *PostgresConnector) GetAllPortfolios(ctx context.Context, limit int, offset int, id int, filterType postgresql.PortfoliosFilterType) ([]models.Portfolio, int, error) {
	portfolios, err := pc.db.GetAllPortfolios(ctx, limit, offset, id, filterType)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get data from database: %w", err)
	}

	rowsAmount, err := pc.db.CountPortfoliosPages(ctx, id, filterType)
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
	return pc.db.GetPortfolioByID(ctx, portfolioID)
}

func (pc *PostgresConnector) CreatePortfolio(ctx context.Context, portfolio models.Portfolio) (int, error) {
	return pc.db.CreatePortfolio(ctx, portfolio)
}

func (pc *PostgresConnector) PatchPortfolio(ctx context.Context, portfolio models.Portfolio) error {
	return pc.db.PatchPortfolio(ctx, portfolio)
}

func (pc *PostgresConnector) DeletePortfolio(ctx context.Context, portfolioID int) error {
	return pc.db.DeletePortfolio(ctx, portfolioID)
}

func (pc *PostgresConnector) CreateCategory(ctx context.Context, name string) (int, error) {
	return pc.db.CreateCategory(ctx, name)
}

func (pc *PostgresConnector) DeleteCategory(ctx context.Context, id int) error {
	return pc.db.DeleteCategory(ctx, id)
}

func (pc *PostgresConnector) GetAllCategories(ctx context.Context, limit int, offset int) ([]models.Category, int, error) {
	categories, err := pc.db.GetAllCategories(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get data from database: %w", err)
	}

	rowsAmount, err := pc.db.CountCategoriesPages(ctx)
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
	crafts, err := pc.db.GetAllCraftsByPortfolioID(ctx, portfolioID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get data from database: %w", err)
	}

	rowsAmount, err := pc.db.CountCraftsPages(ctx, portfolioID, true)
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
	return pc.db.GetCraftByID(ctx, craftID)
}

func (pc *PostgresConnector) CreateCraft(ctx context.Context, portfolioID int, craft models.Craft) (int, error) {
	return pc.db.CreateCraft(ctx, portfolioID, craft)
}

func (pc *PostgresConnector) AddTagToCraft(ctx context.Context, craftID int, tagID int) error {
	return pc.db.AddTagToCraft(ctx, craftID, tagID)
}

func (pc *PostgresConnector) DeleteTagFromCraft(ctx context.Context, craftID int, tagID int) error {
	return pc.db.DeleteTagFromCraft(ctx, craftID, tagID)
}

func (pc *PostgresConnector) PatchCraft(ctx context.Context, craft models.Craft) error {
	return pc.db.PatchCraft(ctx, craft)
}

func (pc *PostgresConnector) DeleteCraft(ctx context.Context, id int) error {
	return pc.db.DeleteCraft(ctx, id)
}

func (pc *PostgresConnector) GetAllCraftsByTagID(ctx context.Context, tagID int, limit int, offset int) ([]models.Craft, int, error) {
	crafts, err := pc.db.GetAllCraftsByTagID(ctx, tagID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get data from database: %w", err)
	}

	rowsAmount, err := pc.db.CountCraftsPages(ctx, tagID, false)
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
	tags, err := pc.db.GetAllTags(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get data from database: %w", err)
	}

	rowsAmount, err := pc.db.CountTagsPages(ctx)
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
	return pc.db.CreateTag(ctx, name)
}

func (pc *PostgresConnector) DeleteTag(ctx context.Context, id int) error {
	return pc.db.DeleteTag(ctx, id)
}

func (pc *PostgresConnector) CreateContent(ctx context.Context, craftID int, content models.Content) (int, error) {
	return pc.db.CreateContent(ctx, craftID, content)
}

func (pc *PostgresConnector) DeleteContent(ctx context.Context, id int) error {
	return pc.db.DeleteContent(ctx, id)
}

func (pc *PostgresConnector) PatchContent(ctx context.Context, content models.Content) error {
	return pc.db.PatchContent(ctx, content)
}
