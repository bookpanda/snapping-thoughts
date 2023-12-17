package twitter

import (
	"context"
	"log"
	"net/http"

	"github.com/dghubble/oauth1"
	"github.com/g8rswimmer/go-twitter/v2"
)

// authorize is not used, but is required by the twitter client
type authorize struct{}

func (a authorize) Add(req *http.Request) {}

type TwitterClient struct {
	client *twitter.Client
}

func NewTwitterClient(
	consumerToken string,
	consumerSecret string,
	userToken string,
	userTokenSecret string,
) *TwitterClient {
	config := oauth1.NewConfig(consumerToken, consumerSecret)
	httpClient := config.Client(oauth1.NoContext, &oauth1.Token{
		Token:       userToken,
		TokenSecret: userTokenSecret,
	})

	client := &twitter.Client{
		Authorizer: authorize{},
		Client:     httpClient,
		Host:       "https://api.twitter.com",
	}

	return &TwitterClient{
		client,
	}
}

func (c *TwitterClient) CreateTweet(text string) (*twitter.CreateTweetResponse, error) {
	req := twitter.CreateTweetRequest{
		Text: text,
	}
	log.Println("Callout to create tweet callout")

	tweetResponse, err := c.client.CreateTweet(context.Background(), req)
	if err != nil {
		log.Panicf("Create tweet error: %v", err)
	}

	return tweetResponse, nil
}
