package api

import (
	"context"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/notifier"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/uptrace/bunrouter"
	"log"
	"net/http"
	"time"

	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/config"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/models"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/storage/postgresql"
)

type Connector interface {
	GetAllPortfolios(ctx context.Context, limit int, offset int, id int, filterType postgresql.PortfoliosFilterType) ([]models.Portfolio, int, error)
	GetPortfolioByID(ctx context.Context, portfolioID int) (*models.Portfolio, error)
	CreatePortfolio(ctx context.Context, portfolio models.Portfolio) (int, error)
	PatchPortfolio(ctx context.Context, portfolio models.Portfolio) error
	DeletePortfolio(ctx context.Context, portfolioID int) error
	CreateCategory(ctx context.Context, name string) (int, error)
	DeleteCategory(ctx context.Context, id int) error
	GetAllCategories(ctx context.Context, limit int, offset int) ([]models.Category, int, error)
	GetAllCraftsByPortfolioID(ctx context.Context, portfolioID int, limit int, offset int) ([]models.Craft, int, error)
	GetCraftByID(ctx context.Context, craftID int) (*models.Craft, error)
	CreateCraft(ctx context.Context, portfolioID int, craft models.Craft) (int, error)
	AddTagToCraft(ctx context.Context, craftID int, tagID int) error
	DeleteTagFromCraft(ctx context.Context, craftID int, tagID int) error
	PatchCraft(ctx context.Context, craft models.Craft) error
	DeleteCraft(ctx context.Context, id int) error
	GetAllCraftsByTagID(ctx context.Context, tagID int, limit int, offset int) ([]models.Craft, int, error)
	GetAllTags(ctx context.Context, limit int, offset int) ([]models.Tag, int, error)
	CreateTag(ctx context.Context, name string) (int, error)
	DeleteTag(ctx context.Context, id int) error
	CreateContent(ctx context.Context, craftID int, content models.Content) (int, error)
	DeleteContent(ctx context.Context, id int) error
	PatchContent(ctx context.Context, content models.Content) error
	NotificationsOn(ctx context.Context, userID int) error
	NotificationsOff(ctx context.Context, userID int) error
	IsNotificationsOn(ctx context.Context, userID int) (bool, error)
}

type Server struct {
	databaseConnector Connector
	notifier          Notifier
	httpServer        *http.Server
}

type Notifier interface {
	Notify(userID int, obj notifier.Object, objID int, change notifier.Change)
}

func NewServer(cfg config.Server, connector Connector, notifier Notifier) *Server {
	s := &Server{
		databaseConnector: connector,
		notifier:          notifier,
	}

	router := bunrouter.New().Compat()
	router.GET("/users/:userID/portfolios", s.getPortfoliosHandler)
	router.GET("/users/:userID/portfolios/:id", s.getPortfolioByIDHandler)
	router.POST("/users/:userID/portfolios", s.postPortfolioHandler)
	router.PATCH("/users/:userID/portfolios/:id", s.patchPortfolioHandler)
	router.DELETE("/users/:userID/portfolios/:id", s.deletePortfolioHandler)

	router.POST("/categories", s.postCategoryHandler)
	router.DELETE("/categories/:id", s.deleteCategoryHandler)
	router.GET("/categories", s.getCategoriesHandler)

	router.GET("/users/:userID/portfolios/:id/crafts", s.getCraftsByPortfolioIDHandler)
	router.GET("/users/:userID/portfolios/:id/crafts/:craftID", s.getCraftHandler)
	router.POST("/users/:userID/portfolios/:id/crafts", s.postCraftHandler)

	router.POST("/users/:userID/portfolios/:id/crafts/:craftID/tags/:tagID", s.postTagPatchCraftHandler)
	router.DELETE("/users/:userID/portfolios/:id/crafts/:craftID/tags/:tagID", s.deleteTagPatchCraftHandler)

	router.PATCH("/users/:userID/portfolios/:id/crafts/:craftID", s.patchCraftHandler)
	router.DELETE("/users/:userID/portfolios/:id/crafts/:craftID", s.deleteCraftHandler)

	router.GET("/tags/:id/crafts", s.getCraftsByTagIDHandler)
	router.GET("/tags", s.getTagsHandler)
	router.POST("/tags", s.postTagHandler)
	router.DELETE("/tags/:id", s.deleteTagHandler)

	router.POST("/users/:userID/portfolios/:id/crafts/:craftID/contents", s.postContentHandler)
	router.DELETE("/users/:userID/portfolios/:id/crafts/:craftID/contents/:contentID", s.deleteContentHandler)
	router.PATCH("/users/:userID/portfolios/:id/crafts/:craftID/contents/:contentID", s.patchContentHandler)

	router.POST("notifications/:userID", s.notificationsOn)
	router.DELETE("notifications/:userID", s.notificationsOff)

	swagHandler := httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json"))
	router.GET("/swagger/*path", swagHandler)

	s.httpServer = &http.Server{
		Addr:         cfg.Listen,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return s
}

func (s *Server) Run() {
	log.Println("server started") // TODO: логгер

	go func() {
		err := s.httpServer.ListenAndServe()
		log.Printf("http server stopped: %s", err.Error()) // TODO: логгер
	}()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
