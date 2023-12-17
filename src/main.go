package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/dghubble/oauth1"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/rs/zerolog/log"

	"github.com/bookpanda/snapping-thoughts/src/client/dynamo"
	"github.com/bookpanda/snapping-thoughts/src/client/tweet"
	seed "github.com/bookpanda/snapping-thoughts/src/seeds"
	"github.com/joho/godotenv"
)

// authorize is not used, but is required by the twitter client
type authorize struct{}

func (a authorize) Add(req *http.Request) {}

func handleArgs(db *dynamo.DynamoDBClient) {
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
	config := oauth1.NewConfig(consumerToken, consumerSecret)
	httpClient := config.Client(oauth1.NoContext, &oauth1.Token{
		Token:       userToken,
		TokenSecret: userTokenSecret,
	})
	twClient := &twitter.Client{
		Authorizer: authorize{},
		Client:     httpClient,
		Host:       "https://api.twitter.com",
	}
	twitterClient := tweet.NewTwitterClient(twClient)

	tableName := os.Getenv("TABLE_NAME")
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(awsSession)
	dynamoClient := dynamo.NewDynamoDBClient(db, tableName)

	handleArgs(dynamoClient)

	item, err := dynamoClient.GetItem()
	if err != nil {
		log.Error().Str("dynamoClient", "get item error").Err(err)
	}
	if item == nil {
		return
	}

	time := time.Now()
	err = dynamoClient.UpdateItem(time, item.Id)
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
	log.Info().Msg(string(enc))
}
