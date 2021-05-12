// gcloud functions deploy password-api --runtime go113 --trigger-http --allow-unauthenticated --source password_cl_ea --entry-point MakeRequest
package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*ATLASURI here*/
var ATLASURI = os.Getenv("MONGODB_URL_START") + os.Getenv("MONGODB_CRYPTO_PASSWORD") + os.Getenv("MONGODB_URL_END")

// client is used to make HTTP requests with a 10 second timeout.
// http.Clients should be reused instead of created as needed.
var client = &http.Client{
	Timeout: 10 * time.Second,
}

var ValidKeys = []string{"password"}

type Result struct {
	Password string `json:"password"`
	Answer   string `json:"answer"`
}

type RequestData struct {
	JobRunID string  `json:"jobRunID`
	Data     bson.M  `json:"data"`
	Result   float64 `json:"result"`
	Status   string  `json:"status"`
}

type ChainlinkResult struct {
	JobRunID string  `json:"jobRunID`
	Data     Result  `json:"data"`
	Result   float64 `json:"result"`
	Status   string  `json:"status"`
	Error    string  `json:"error"`
}

func MakeRequest(w http.ResponseWriter, r *http.Request) {
	// get input parameters and query parameters
	parameters := make(map[string][]string)
	parameters["jobRunID"] = []string{""}
	reqBody, err := ioutil.ReadAll(r.Body)
	var data RequestData
	json.Unmarshal(reqBody, &data)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range data.Data {
		parameters[key] = []string{value.(string)}
	}

	for k, v := range r.URL.Query() {
		parameters[k] = v
	}
	log.Printf("%v", parameters)
	filter := getFilterFromParamters(parameters)
	log.Printf("%v", filter)

	if _, ok := filter["password"]; !ok {
		fail_result := ChainlinkResult{Status: "404", Error: "You need 'password'"}
		if err := json.NewEncoder(w).Encode(&fail_result); err != nil {
			logrus.Errorf("Failed to encode response: %v", err)
		}
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(ATLASURI))
	databases, _ := client.ListDatabaseNames(ctx, bson.M{})
	log.Printf("%v", databases)
	collections, _ := client.Database("DATABASES").ListCollectionNames(ctx, bson.M{})
	log.Printf("%v", collections)

	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("COLLECTION"))

	result := getResultFromDatabase(collection, filter)
	log.Printf("%v", result)
	// Decode the result
	w.Header().Set("Content-Type", "application/json")
	result["result"] = result["answer"]
	chainlinkReturn := bson.M{"jobRunID": parameters["jobRunID"][0], "data": result, "status": "200", "result": result["answer"]}
	log.Printf("%v", chainlinkReturn)
	if err := json.NewEncoder(w).Encode(&chainlinkReturn); err != nil {
		logrus.Errorf("Failed to encode response: %v", err)
	}
	client.Disconnect(ctx)
}

func stringInSlice(stringToCheck string, listOfWords []string) bool {
	for _, word := range listOfWords {
		if word == stringToCheck {
			return true
		}
	}
	return false
}

func getFilterFromParamters(parameters map[string][]string) bson.M {
	filter := bson.M{}
	for key, values := range parameters {
		if stringInSlice(key, ValidKeys) {
			filter[key] = values[0]
		}
	}
	return filter
}

func getResultFromDatabase(collection *mongo.Collection, filter bson.M) map[string]interface{} {
	log.Printf("%v", filter)
	var result map[string]interface{}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Printf("%v", err)
	}
	return result
}

func addParametersFromPost(parameters map[string][]string, data interface{}) map[string][]string {
	for k, value := range data.(map[string][]string) {
		parameters[k] = value
	}
	return parameters
}

// curl -X POST -H "content-type:application/json" "https://us-central1-alpha-chain-api.cloudfunctions.net/test-api?to_symbol=USD&from_symbol=LINK"
// curl -X POST -H "content-type:application/json" "https://us-central1-alpha-chain-api.cloudfunctions.net/test-api?to_symbol=ETH&from_symbol=XDR"
// $ curl -X POST -H "content-type:application/json" "https://us-central1-alpha-chain-api.cloudfunctions.net/develop-test-api" --data '{"data":{"from_symbol":"XDR", "to_symbol":"ETH", "chainlink_node":"true"}}'
//  curl -X GET "https://us-central1-alpha-chain-api.cloudfunctions.net/data-query?to_symbol=USD&from_symbol=XAG"
