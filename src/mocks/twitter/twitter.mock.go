package twitter

import (
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) CreateTweet(text string) (res *twitter.CreateTweetResponse, err error) {
	args := c.Called(text)

	if args.Get(0) != nil {
		res = args.Get(0).(*twitter.CreateTweetResponse)
	}

	return res, args.Error(1)
}
