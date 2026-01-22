package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func DBInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Unable to find .env")
	}

	MongoDb_URI := os.Getenv("MONGO_URI")
	if MongoDb_URI == "" {
		log.Fatal("MongoDB URI not set")
	}
	fmt.Println("MonogoDB URI", MongoDb_URI)

	clientOption := options.Client().ApplyURI(MongoDb_URI)
	client, err := mongo.Connect(clientOption)
	if err != nil {
		return nil
	}

	return client

}

var client *mongo.Client = DBInstance()

func OpenCollection(collectionName string) *mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("MongoDB URI not set")
	}
	databaseName := os.Getenv("DATABASE_NAME")
	fmt.Println("Database Name :", databaseName)

	collection := client.Database(databaseName).Collection(collectionName)
	if collection == nil {
		return  nil
	}
	return collection
}
