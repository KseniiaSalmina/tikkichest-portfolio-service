package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/config"
)

type DB struct {
	db *pgxpool.Pool
}

func NewDB(ctx context.Context, cfg config.Postgres) (*DB, error) {
	connstr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := pgxpool.New(ctx, connstr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{
		db: db,
	}, nil
}

func (db *DB) Close() {
	db.db.Close()
}

func GetProfileIDByPortfolio() {} // TODO: для перехода на профиль автора портфолио
func GetProfileIDByCraft()     {} // TODO: для перехода на профиль автора крафта
func GetPortfolioIDByCraft()   {} // TODO: для перехода на портфолио по найденному крафту
// TODO: модифицировать методы обновления объектов для обновления не всех параметров за раз
