package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/mock"
)

type DatabaseMock struct {
	mock.Mock
}

func (c *DatabaseMock) PutItem(in *dynamodb.PutItemInput) (res *dynamodb.PutItemOutput, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*dynamodb.PutItemOutput)
	}

	return res, args.Error(1)
}

func (c *DatabaseMock) Scan(in *dynamodb.ScanInput) (res *dynamodb.ScanOutput, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*dynamodb.ScanOutput)
	}

	return res, args.Error(1)
}

func (c *DatabaseMock) GetItem(in *dynamodb.GetItemInput) (res *dynamodb.GetItemOutput, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*dynamodb.GetItemOutput)
	}

	return res, args.Error(1)
}

func (c *DatabaseMock) UpdateItem(in *dynamodb.UpdateItemInput) (res *dynamodb.UpdateItemOutput, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*dynamodb.UpdateItemOutput)
	}

	return res, args.Error(1)
}
