package tweet

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/g8rswimmer/go-twitter/v2"
)

type TwitterClient struct {
	client Client
}

type Client interface {
	CreateTweet(ctx context.Context, tweet twitter.CreateTweetRequest) (*twitter.CreateTweetResponse, error)
}

func NewTwitterClient(
	client Client,
) *TwitterClient {
	return &TwitterClient{
		client,
	}
}

func (c *TwitterClient) CreateTweet(text string) (*twitter.CreateTweetResponse, error) {
	req := &twitter.CreateTweetRequest{
		Text: text,
	}

	tweetResponse, err := c.client.CreateTweet(context.Background(), *req)
	if err != nil {
		log.Fatal().Str("twitterClient", "Create tweet error").Err(err)
		return nil, err
	}

	return tweetResponse, nil
}
