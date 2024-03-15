package api

import (
	"context"
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
}

type Server struct {
	databaseConnector Connector
	httpServer        *http.Server
}

func NewServer(cfg config.Server, connector Connector) *Server {
	s := &Server{
		databaseConnector: connector,
	}

	router := bunrouter.New().Compat()
	router.GET("/portfolios", s.getPortfoliosHandler)
	router.GET("/portfolios/:id", s.getPortfolioByIDHandler)
	router.POST("/portfolios", s.postPortfolioHandler)
	router.PATCH("/portfolios/:id", s.patchPortfolioHandler)
	router.DELETE("/portfolios/:id", s.deletePortfolioHandler)
	router.POST("/categories", s.postCategoryHandler)
	router.DELETE("/categories/:id", s.deleteCategoryHandler)
	router.GET("/categories", s.getCategoriesHandler)
	router.GET("/portfolios/:id/crafts", s.getCraftsByPortfolioIDHandler)
	router.GET("/portfolios/:id/crafts/:craftID", s.getCraftHandler)
	router.POST("/portfolios/:id/crafts", s.postCraftHandler)
	router.POST("/portfolios/:id/crafts/:craftID/tags/:tagID", s.postTagPatchCraftHandler)     
	router.DELETE("/portfolios/:id/crafts/:craftID/tags/:tagID", s.deleteTagPatchCraftHandler) 
	router.PATCH("/portfolios/:id/crafts/:craftID", s.patchCraftHandler)                       
	router.DELETE("/portfolios/:id/crafts/:craftID", s.deleteCraftHandler)
	router.GET("/tags/:id/crafts", s.getCraftsByTagIDHandler) 
	router.GET("/tags", s.getTagsHandler)
	router.POST("/tags", s.postTagHandler)
	router.DELETE("/tags/:id", s.deleteTagHandler)
	router.POST("/portfolios/:id/crafts/:craftID/contents", s.postContentHandler)
	router.DELETE("/portfolios/:id/crafts/:craftID/contents/:contentID", s.deleteContentHandler)
	router.PATCH("/portfolios/:id/crafts/:craftID/contents/:contentID", s.patchContentHandler) 

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
