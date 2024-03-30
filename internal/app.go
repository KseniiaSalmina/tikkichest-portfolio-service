package app

import (
	"context"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/notifier"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/notifier/sender/kafka"
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
	notifier *notifier.Notifier
	sender *kafka.ProducerManager
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
	if err := a.initSender(); err != nil {
		return err
	}
	a.initNotifier()

	//init controllers
	a.initServer()

	return nil
}

func (a *Application) initDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

func (a *Application) initSender() error {
	sender, err := kafka.NewProducerManager(a.cfg.Kafka)
	if err != nil {
		log.Println(err) // TODO: logger
		return err
	}

	a.sender = sender
	return nil
}

func (a *Application) initNotifier() {
	a.notifier = notifier.NewNotifier(a.sender)
}

func (a *Application) initServer() {
	s := api.NewServer(a.cfg.Server, a.dbConnector, a.notifier)

	a.server = s
}

func (a *Application) Run() {
	defer a.stop()

	a.sender.Run(a.closeCtx)
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

	if err := a.sender.Shutdown(); err != nil {
		log.Print(err) // TODO: logger
	}

	a.db.Close()
	log.Print("database closed") // TODO: logger
}

func (a *Application) readyToShutdown() {
	ctx, closeCtx := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	a.closeCtx = ctx
	a.closeCtxFunc = closeCtx
}
