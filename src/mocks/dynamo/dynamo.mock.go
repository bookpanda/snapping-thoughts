package dynamo

import (
	"github.com/bookpanda/snapping-thoughts/src/model/item"
	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) CreateItem(in *item.Item) error {
	args := c.Called(in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*item.Item)
	}

	return args.Error(1)
}

func (c *ClientMock) GetItem() (res *item.Item, err error) {
	args := c.Called()

	if args.Get(0) != nil {
		res = args.Get(0).(*item.Item)
	}

	return res, args.Error(1)
}

func (c *ClientMock) GetItemWithId(id string) (res *item.Item, err error) {
	args := c.Called(id)

	if args.Get(0) != nil {
		res = args.Get(0).(*item.Item)
	}

	return res, args.Error(1)
}

func (c *ClientMock) UpdateItem(id string) error {
	args := c.Called(id)

	return args.Error(0)
}
