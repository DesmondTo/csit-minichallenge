package flight

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/DesmondTo/minichallenge/internal/model"
)

func Get(collection *mongo.Collection, filters bson.D, sortOptions bson.D, projectOptions bson.D) ([]model.Flight, error) {
	opts := options.Find().SetSort(sortOptions).SetProjection(projectOptions)
	cursor, err := collection.Find(context.Background(), filters, opts)
	if err != nil {
		return nil, err
	}

	var results []model.Flight
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
