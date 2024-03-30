package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/config"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/models"
)

type DB struct {
	*mongo.Collection
}

func NewDB(cfg config.Mongo) (*mongo.Collection, error) {
	var connectStr string
	var connectOpt options.ClientOptions

	if cfg.User == "" && cfg.Password == "" {
		connectStr = fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
	} else {
		connectStr = fmt.Sprintf("mongodb://%s:%s@%s:%d", cfg.User, cfg.Password, cfg.Host, cfg.Port)
		if cfg.AuthenticationDatabase == "" {
			cfg.AuthenticationDatabase = cfg.Database
		}
		connectOpt.SetAuth(options.Credential{Username: cfg.User, Password: cfg.Password, AuthSource: cfg.AuthenticationDatabase})
	}

	client, err := mongo.Connect(context.Background(), connectOpt.ApplyURI(connectStr))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	timeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() //TODO????

	if err = client.Ping(timeout, nil); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := client.Database(cfg.Database).Collection(cfg.Collection)

	if !cfg.HaveIndexes {
		models := []mongo.IndexModel{
			{Keys: bson.D{{"profile_id", 1}}, Options: options.Index().SetName("profile_id_idx")},
			{Keys: bson.D{{"category", 1}}, Options: options.Index().SetName("category_idx")},
		}

		timeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel() //TODO????

		if _, err = db.Indexes().CreateMany(timeout, models); err != nil {
			return nil, fmt.Errorf("failed to create profile id index: %w", err)
		}
	}

	return db, nil
}

// ByProfileID returns filter by profile ID
func ByProfileID(id int) bson.D {
	return constructor("profile_id", id)
}

// ByCategory returns filter by category
func ByCategory(category string) bson.D {
	return constructor("category", category)
}

func constructor(key string, value interface{}) bson.D {
	return bson.D{{key, value}}
}

func CreatePortfolio(ctx context.Context, db *mongo.Collection, portfolio models.Portfolio) error {
	if _, err := db.InsertOne(ctx, portfolio); err != nil {
		return fmt.Errorf("failed to create portfolio: %w", err)
	}

	return nil
}

func DeletePortfolio(ctx context.Context, db *mongo.Collection, portfolioID string) error {
	if _, err := db.DeleteOne(ctx, constructor("_id", portfolioID)); err != nil {
		return fmt.Errorf("failed to delete portfolio: %w", err)
	}

	return nil
}

func GetAllPortfolios(ctx context.Context, db *mongo.Collection, filter ...bson.D) ([]models.Portfolio, error) {
	if filter == nil {
		filter = []bson.D{{}}
	}

	cursor, err := db.Find(ctx, filter[0], options.Find().SetProjection(bson.D{{"crafts", 0}}))

	var portfolios []models.Portfolio
	if err = cursor.All(context.Background(), &portfolios); err != nil {
		return nil, fmt.Errorf("failed to get portfolios: %w", err)
	}

	return portfolios, nil
}

func GetPortfolio(ctx context.Context, db *mongo.Collection, portfolioID string) (*models.Portfolio, error) {
	filter := constructor("_id", portfolioID)

	var portfolio models.Portfolio
	if err := db.FindOne(ctx, filter, options.FindOne().SetProjection(bson.D{{"crafts", 0}})).Decode(&portfolio); err != nil {
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}

	return &portfolio, nil
}

func PatchPortfolio(ctx context.Context, db *mongo.Collection, portfolio models.Portfolio) error {
	filter := constructor("_id", portfolio.ID)
	update := bson.D{{"$set", bson.D{{"name", portfolio.Name}, {"category", portfolio.Category}, {"description", portfolio.Description}}}}

	var oldPortfolio models.Portfolio
	if err := db.FindOneAndUpdate(ctx, filter, update).Decode(&oldPortfolio); err != nil {
		return fmt.Errorf("failed to update portfolio: %w", err)
	}

	return nil
}

func CreateCraft(ctx context.Context, db *mongo.Collection, portfolio models.Portfolio) error {
	filter := constructor("_id", portfolio.ID)
	update := constructor("$push", constructor("crafts", portfolio.Crafts[0]))

	var oldPortfolio models.Portfolio
	if err := db.FindOneAndUpdate(ctx, filter, update).Decode(&oldPortfolio); err != nil {
		return fmt.Errorf("failed to create a new craft: %w", err)
	}

	return nil
}

func DeleteCraft(ctx context.Context, db *mongo.Collection, portfolio models.Portfolio) error {
	filter := constructor("_id", portfolio.ID)
	update := constructor("$pull", constructor("crafts", constructor("_id", portfolio.Crafts[0].ID)))

	var oldPortfolio models.Portfolio
	if err := db.FindOneAndUpdate(ctx, filter, update).Decode(&oldPortfolio); err != nil {
		return fmt.Errorf("failed to delete craft: %w", err)
	}
	return nil
}

func PatchCraft(ctx context.Context, db *mongo.Collection, portfolio models.Portfolio) error {
	filter := bson.D{{"crafts._id", portfolio.Crafts[0].ID}}
	update := constructor("$set", constructor("crafts.name", portfolio.Crafts[0].Name))

	var oldCraft models.Craft
	if err := db.FindOneAndUpdate(ctx, filter, update).Decode(&oldCraft); err != nil {
		return fmt.Errorf("failed to update craft: %w", err)
	}
	//TODO
	return nil
}

func GetCraft(ctx context.Context, db *mongo.Collection, portfolio models.Portfolio) (*models.Craft, error) {
	filter := constructor("_id", portfolio.ID)

	var craft models.Craft
	if err := db.FindOne(ctx, filter).Decode(&craft); err != nil {
		return nil, fmt.Errorf("failed to get craft: %w", err)
	}

	return &craft, nil
}

func GetAllCrafts(ctx context.Context, db *mongo.Collection, portfolioID int) (*models.Portfolio, error) {
	filter := constructor("_id", portfolioID)

	var portfolio models.Portfolio
	if err := db.FindOne(ctx, filter).Decode(&portfolio); err != nil {
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}

	for i := range portfolio.Crafts {
		portfolio.Crafts[i].Contents = portfolio.Crafts[i].Contents[:1]
	}

	return &portfolio, nil
}

func CreateContent(ctx context.Context, db *mongo.Collection, craft models.Craft) error {
	filter := constructor("craft_id", craft.ID)
	update := bson.D{{"$push", bson.D{{"contents", craft.Contents[0]}}}}

	var oldCraft models.Craft
	if err := db.FindOneAndUpdate(ctx, filter, update).Decode(&oldCraft); err != nil {
		return fmt.Errorf("failed to create a new content: %w", err)
	}

	return nil
}

func DeleteContent(ctx context.Context, db *mongo.Collection, craft models.Craft) error {
	filter := constructor("craft_id", craft.ID)
	update := bson.D{{"$pull", constructor("content_id", craft.Contents[0].ID)}}

	var oldCraft models.Craft
	if err := db.FindOneAndUpdate(ctx, filter, update).Decode(&oldCraft); err != nil {
		return fmt.Errorf("failed to delete content: %w", err)
	}

	return nil
}

func PatchContent(ctx context.Context, db *mongo.Collection, content models.Content) error {
	filter := constructor("content_id", content.ID)
	update := bson.D{{"$set", bson.D{{"content_description", content.Description}}}}

	var oldContent models.Content
	if err := db.FindOneAndUpdate(ctx, filter, update).Decode(&oldContent); err != nil {
		return fmt.Errorf("failed to update content: %w", err)
	}

	return nil
}
