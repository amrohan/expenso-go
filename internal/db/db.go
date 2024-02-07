package db

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientInstance *mongo.Client

var mongoOnce sync.Once

var clientInstanceError error

type Collection string

const (
	TransactionCollection Collection = "transactions"
	CategoryCollection    Collection = "categories"
	AccountCollection     Collection = "accounts"
	UserCollection        Collection = "users"
)

const (
	Database = "budget"
)

func GetMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {

		err := godotenv.Load()

		if err != nil {
			log.Printf("Error loading .env file")
		}

		MongoURl := os.Getenv("MONGO_URL")

		if MongoURl == "" {
			panic("Mongo URL is empty")
		}

		clientOptions := options.Client().ApplyURI(MongoURl)

		client, err := mongo.Connect(context.TODO(), clientOptions)

		clientInstance = client

		clientInstanceError = err
	})

	return clientInstance, clientInstanceError
}
