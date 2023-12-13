package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bookpanda/snapping-thoughts/src/client"
	"github.com/joho/godotenv"
)

type authorize struct{}

func (a authorize) Add(req *http.Request) {}

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

	text := flag.String("text", "hello3", "twitter text")
	flag.Parse()

	tweetResponse, err := twitterClient.CreateTweet(*text)
	if err != nil {
		log.Panicf("create tweet error: %v", err)
	}

	enc, err := json.MarshalIndent(tweetResponse, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(enc))
}
