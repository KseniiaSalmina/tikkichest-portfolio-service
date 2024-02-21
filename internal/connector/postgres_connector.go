package connector

import (
	"context"
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

func (pc *PostgresConnector) GetAllPortfolios(ctx context.Context, limit int, offset int, filter postgresql.PortfoliosFilter) ([]models.Portfolio, error) {
	return postgresql.GetAllPortfolios(ctx, pc.db, limit, offset, filter)
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

func (pc *PostgresConnector) GetAllCategories(ctx context.Context, limit int, offset int) ([]models.Category, error) {
	return postgresql.GetAllCategories(ctx, pc.db, limit, offset)
}

func (pc *PostgresConnector) GetAllCraftsByPortfolioID(ctx context.Context, portfolioID int, limit int, offset int) ([]models.Craft, error) {
	return postgresql.GetAllCraftsByPortfolioID(ctx, pc.db, portfolioID, limit, offset)
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

func (pc *PostgresConnector) GetAllCraftsByTagID(ctx context.Context, tagID int, limit int, offset int) ([]models.Craft, error) {
	return postgresql.GetAllCraftsByTagID(ctx, pc.db, tagID, limit, offset)
}

func (pc *PostgresConnector) GetAllTags(ctx context.Context, limit int, offset int) ([]models.Tag, error) {
	return postgresql.GetAllTags(ctx, pc.db, limit, offset)
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
