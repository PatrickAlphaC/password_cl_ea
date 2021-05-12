package http

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestConnect(t *testing.T) {
	ATLAS_URI := os.Getenv("MONGODB_URL_START") + os.Getenv("MONGODB_CRYPTO_PASSWORD") + os.Getenv("MONGODB_URL_END")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := mongo.Connect(ctx, options.Client().ApplyURI(ATLAS_URI))
	if err != nil {
		t.Errorf("Error connecting: %v", err)
	}
}

func TestGetFilter(t *testing.T) {
	parameters := make(map[string][]string)
	parameters["password"] = []string{"test_password"}
	filter := getFilterFromParamters(parameters)
	if len(filter) != 1 {
		t.Errorf(("Error getting filter"))
	}
	t.Logf("%v", filter)
}

func TestGetEntryFromFilter(t *testing.T) {
	ATLAS_URI := os.Getenv("MONGODB_URL_START") + os.Getenv("MONGODB_CRYPTO_PASSWORD") + os.Getenv("MONGODB_URL_END")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(ATLAS_URI))
	if err != nil {
		t.Errorf("Error connecting: %v", err)
	}
	collections, _ := client.Database(os.Getenv("DATABASE")).ListCollectionNames(ctx, bson.M{})
	log.Printf("%v", collections)

	collection := client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("COLLECTION"))
	filter := bson.M{"password": "test_password"}
	log.Printf("%v", filter)
	result := getResultFromDatabase(collection, filter)
	t.Logf("%v", result)
	t.Logf("%v", result)
	if result["answer"].(string) != "test_answer" {
		t.Errorf("Error getting the right stuff from the DB")
	}
	t.Logf("%v", result)
}

func TestGettingURLParameters(t *testing.T) {
	ATLAS_URI := os.Getenv("MONGODB_URL_START") + os.Getenv("MONGODB_CRYPTO_PASSWORD") + os.Getenv("MONGODB_URL_END")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(ATLAS_URI))
	if err != nil {
		t.Errorf("Error connecting: %v", err)
	}
	collection := client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("COLLECTION"))
	filter := bson.M{"password": "test_password"}
	result := getResultFromDatabase(collection, filter)
	if result["answer"].(string) != "test_answer" {
		t.Errorf("Error getting the right stuff from the DB")
	}
	t.Logf("%v", result)
}
