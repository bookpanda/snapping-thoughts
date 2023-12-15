package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bookpanda/snapping-thoughts/src/client"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnv()
	consumerToken := os.Getenv("CONSUMER_API_KEY")
	consumerSecret := os.Getenv("CONSUMER_API_SECRET")
	userToken := os.Getenv("ACCESS_TOKEN")
	userTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	twitterClient := client.NewTwitterClient(consumerToken, consumerSecret, userToken, userTokenSecret)

	tableName := os.Getenv("TABLE_NAME")
	dynamoClient := client.NewDynamoDBClient(tableName)

	item, err := dynamoClient.GetItem()
	if err != nil {
		log.Panicf("get item error: %v", err)
	}
	if item == nil {
		return
	}

	err = dynamoClient.UpdateItem(item.Id)
	if err != nil {
		log.Panicf("update item error: %v", err)
	}

	tweetResponse, err := twitterClient.CreateTweet(item.Message)
	if err != nil {
		log.Panicf("create tweet error: %v", err)
	}

	enc, err := json.MarshalIndent(tweetResponse, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(enc))
}
