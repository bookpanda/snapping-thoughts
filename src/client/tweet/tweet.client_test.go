package tweet

import (
	"errors"
	"testing"

	mock "github.com/bookpanda/snapping-thoughts/src/mocks/tweet"
	"github.com/bxcodec/faker/v3"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TwitterClientTest struct {
	suite.Suite
	CreateTweetInput *twitter.CreateTweetRequest
	Text             string
}

func TestTwitterClient(t *testing.T) {
	suite.Run(t, new(TwitterClientTest))
}

func (t *TwitterClientTest) SetupTest() {
	t.Text = faker.Word()
	t.CreateTweetInput = &twitter.CreateTweetRequest{
		Text: t.Text,
	}
}

func (t *TwitterClientTest) TestCreateItemSuccess() {
	output := &twitter.CreateTweetResponse{
		Tweet: &twitter.CreateTweetData{
			Text: t.Text,
		},
	}

	mockClient := mock.ClientMock{}
	mockClient.On("CreateTweet", *t.CreateTweetInput).Return(output, nil)

	client := NewTwitterClient(&mockClient)
	res, err := client.CreateTweet(t.Text)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), output.Tweet.Text, res.Tweet.Text)
}

func (t *TwitterClientTest) TestCreateItemInternalErr() {
	output := &twitter.CreateTweetResponse{}

	mockClient := mock.ClientMock{}
	mockClient.On("CreateTweet", *t.CreateTweetInput).Return(output, errors.New("something wrong"))

	client := NewTwitterClient(&mockClient)
	_, err := client.CreateTweet(t.Text)

	assert.Error(t.T(), err)
}
