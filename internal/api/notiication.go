package api

import (
	"context"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/notifier"
	"log"
)

func (s *Server) Notify(ctx context.Context, userID int, obj notifier.Object, objID int, change notifier.Change) {
	on, err := s.databaseConnector.IsNotificationsOn(ctx, userID)
	if err != nil {
		log.Println(err) //TODO logger
	}
	if on {
		s.notifier.Notify(userID, obj, objID, change)
	}
}