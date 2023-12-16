package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/bookpanda/snapping-thoughts/src/client"
	seed "github.com/bookpanda/snapping-thoughts/src/seeds"
	"github.com/joho/godotenv"
)

func handleArgs(db *client.DynamoDBClient) {
	flag.Parse()
	args := flag.Args()

	if len(args) >= 1 {
		switch args[0] {
		case "seed":
			err := seed.Execute(db, args[1:]...)
			if err != nil {
				log.Fatal().
					Str("service", "seeder").
					Msg("Not found seed")
			}
			os.Exit(0)
		}
	}
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
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
	handleArgs(dynamoClient)

	item, err := dynamoClient.GetItem()
	if err != nil {
		log.Error().Str("dynamoClient", "get item error").Err(err)
	}
	if item == nil {
		return
	}

	err = dynamoClient.UpdateItem(item.Id)
	if err != nil {
		log.Error().Str("dynamoClient", "update item error").Err(err)
	}

	tweetResponse, err := twitterClient.CreateTweet(item.Message)
	if err != nil {
		log.Error().Str("twitterClient", "create tweet error").Err(err)
	}

	enc, err := json.MarshalIndent(tweetResponse, "", "    ")
	if err != nil {
		log.Error().Str("twitterClient", "marshal response error").Err(err)
	}
	fmt.Println(string(enc))
}
