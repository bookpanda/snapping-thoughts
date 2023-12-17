package dynamo

import (
	"github.com/bookpanda/snapping-thoughts/src/model/dynamo"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (c *RepositoryMock) CreateItem(in *dynamo.Item) error {
	args := c.Called(in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*dynamo.Item)
	}

	return args.Error(1)
}

func (c *RepositoryMock) GetItem() (res *dynamo.Item, err error) {
	args := c.Called()

	if args.Get(0) != nil {
		res = args.Get(0).(*dynamo.Item)
	}

	return res, args.Error(1)
}

func (c *RepositoryMock) GetItemWithId(id string) (res *dynamo.Item, err error) {
	args := c.Called(id)

	if args.Get(0) != nil {
		res = args.Get(0).(*dynamo.Item)
	}

	return res, args.Error(1)
}

func (c *RepositoryMock) UpdateItem(id string) error {
	args := c.Called(id)

	return args.Error(0)
}
