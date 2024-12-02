package common

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitializeMongoDb() *mongo.Client {

	mongoDbUrl := os.Getenv("MONGODB_URL")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDbUrl))

	FailOnError(err, "Failed to connect to mongo db")

	return client

}
