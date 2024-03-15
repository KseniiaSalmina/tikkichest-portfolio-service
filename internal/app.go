package app

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/api"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/config"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/connector"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/storage/postgresql"
)

type Application struct {
	cfg          config.Application
	db           *postgresql.DB
	dbConnector  *connector.PostgresConnector
	server       *api.Server
	closeCtx     context.Context
	closeCtxFunc context.CancelFunc
}

func NewApplication(cfg config.Application) (*Application, error) {
	app := Application{
		cfg: cfg,
	}

	if err := app.bootstrap(); err != nil {
		return nil, err
	}

	app.readyToShutdown()

	return &app, nil
}

func (a *Application) bootstrap() error {
	//init dependencies
	if err := a.initDatabase(); err != nil {
		return err
	}

	//init services
	a.initConnector()

	//init controllers
	if err := a.initServer(); err != nil {
		return err
	}

	return nil
}

func (a *Application) initDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second) // TODO: correct time
	defer cancel()

	db, err := postgresql.NewDB(ctx, a.cfg.Storage.Postgres)
	if err != nil {
		log.Println(err) // TODO: logger
		return err
	}

	a.db = db
	log.Println("successful connection to database") // TODO: logger
	return nil
}

func (a *Application) initConnector() {
	a.dbConnector = connector.NewPostgresConnector(a.db)
}

func (a *Application) initServer() error {
	s := api.NewServer(a.cfg.Server, a.dbConnector)

	a.server = s
	return nil
}

func (a *Application) Run() {
	defer a.stop()

	a.server.Run()

	<-a.closeCtx.Done()
	a.closeCtxFunc()
}

func (a *Application) stop() {
	if err := a.server.Shutdown(); err != nil {
		log.Printf("incorrect closing of server: %s", err.Error()) // TODO: logger
	} else {
		log.Print("server closed") // TODO: logger
	}

	a.db.Close()
	log.Print("database closed") // TODO: logger
}

func (a *Application) readyToShutdown() {
	ctx, closeCtx := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	a.closeCtx = ctx
	a.closeCtxFunc = closeCtx
}
