package dynamo

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	mock "github.com/bookpanda/snapping-thoughts/src/mocks/dynamo"
	"github.com/bookpanda/snapping-thoughts/src/model/dynamo"
	"github.com/bxcodec/faker/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DynamoClientTest struct {
	suite.Suite
	Item            *dynamo.Item
	PutItemInput    *dynamodb.PutItemInput
	ScanItemInput   *dynamodb.ScanInput
	GetItemInput    *dynamodb.GetItemInput
	GetItemId       string
	UpdateItemInput *dynamodb.UpdateItemInput
	UpdateItemId    string
	TableName       string
}

func TestDynamoClient(t *testing.T) {
	suite.Run(t, new(DynamoClientTest))
}

func (t *DynamoClientTest) SetupTest() {
	t.Item = &dynamo.Item{
		Id:      faker.UUIDHyphenated(),
		IsUsed:  faker.Word(),
		Message: faker.Word(),
	}

	av, err := dynamodbattribute.MarshalMap(t.Item)
	if err != nil {
		log.Fatal().Str("twitterClient test", "Got error marshalling new item").Err(err)
	}

	t.PutItemInput = &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(t.TableName),
	}

	proj := expression.NamesList(expression.Name("Message"), expression.Name("Id"))
	filt := expression.Name("IsUsed").Equal(expression.Value("no"))
	expr, err := expression.NewBuilder().WithProjection(proj).WithFilter(filt).Build()
	if err != nil {
		log.Fatal().Str("twitterClient", "Got error building expression").Err(err)
	}

	t.ScanItemInput = &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(t.TableName),
		Limit:                     aws.Int64(1),
	}

	t.GetItemId = faker.UUIDHyphenated()
	t.GetItemInput = &dynamodb.GetItemInput{
		TableName: aws.String(t.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(t.GetItemId),
			},
		},
	}

	t.UpdateItemId = faker.UUIDHyphenated()
	key := map[string]*dynamodb.AttributeValue{
		"Id": {
			S: aws.String(t.UpdateItemId),
		},
	}
	t.UpdateItemInput = &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":u": {
				S: aws.String(time.Now().String()),
			},
		},
		TableName:        aws.String(t.TableName),
		Key:              key,
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set IsUsed = :u"),
	}

	t.TableName = faker.Word()
}

func (t *DynamoClientTest) TestCreateItemSuccess() {
	output := &dynamodb.PutItemOutput{}

	db := mock.DatabaseMock{}
	db.On("PutItem", &t.PutItemInput).Return(output, nil)

	client := NewDynamoDBClient(&db, t.TableName)
	err := client.CreateItem(t.Item)

	assert.Nil(t.T(), err)
}

func (t *DynamoClientTest) TestCreateItemInternalErr() {
	output := &dynamodb.PutItemOutput{}

	db := mock.DatabaseMock{}
	db.On("PutItem", &t.PutItemInput).Return(output, errors.New("something wrong"))

	client := NewDynamoDBClient(&db, t.TableName)
	err := client.CreateItem(t.Item)

	assert.Error(t.T(), err)
}

func (t *DynamoClientTest) TestGetItemSuccess() {
	want := t.Item
	output := &dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"Id": {
					S: aws.String(t.Item.Id),
				},
				"IsUsed": {
					S: aws.String(t.Item.IsUsed),
				},
				"Message": {
					S: aws.String(t.Item.Message),
				},
			},
		},
	}

	db := mock.DatabaseMock{}
	db.On("Scan", &t.ScanItemInput).Return(output, nil)

	client := NewDynamoDBClient(&db, t.TableName)
	res, err := client.GetItem()

	assert.Equal(t.T(), want, res)
	assert.Nil(t.T(), err)
}

func (t *DynamoClientTest) TestGetItemInternalErr() {
	output := &dynamodb.ScanOutput{}

	db := mock.DatabaseMock{}
	db.On("Scan", &t.ScanItemInput).Return(output, errors.New("something wrong"))

	client := NewDynamoDBClient(&db, t.TableName)
	_, err := client.GetItem()

	assert.Error(t.T(), err)
}

func (t *DynamoClientTest) TestGetItemWithIdSuccess() {
	want := t.Item
	output := &dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(t.Item.Id),
			},
			"IsUsed": {
				S: aws.String(t.Item.IsUsed),
			},
			"Message": {
				S: aws.String(t.Item.Message),
			},
		},
	}

	db := mock.DatabaseMock{}
	db.On("GetItem", &t.GetItemInput).Return(output, nil)

	client := NewDynamoDBClient(&db, t.TableName)
	res, err := client.GetItemWithId(t.GetItemId)

	assert.Equal(t.T(), want, res)
	assert.Nil(t.T(), err)
}

func (t *DynamoClientTest) TestGetItemWithIdInternalErr() {
	output := &dynamodb.GetItemOutput{}

	db := mock.DatabaseMock{}
	db.On("GetItem", &t.GetItemInput).Return(output, errors.New("something wrong"))

	client := NewDynamoDBClient(&db, t.TableName)
	_, err := client.GetItemWithId(t.GetItemId)

	assert.Error(t.T(), err)
}

func (t *DynamoClientTest) TestUpdateItemSuccess() {
	output := &dynamodb.UpdateItemOutput{}

	db := mock.DatabaseMock{}
	db.On("UpdateItem", &t.UpdateItemInput).Return(output, nil)

	client := NewDynamoDBClient(&db, t.TableName)
	err := client.UpdateItem(t.UpdateItemId)

	assert.Nil(t.T(), err)
}

func (t *DynamoClientTest) TestUpdateItemInternalErr() {
	output := &dynamodb.UpdateItemOutput{}

	db := mock.DatabaseMock{}
	db.On("UpdateItem", &t.UpdateItemInput).Return(output, errors.New("something wrong"))

	client := NewDynamoDBClient(&db, t.TableName)
	err := client.UpdateItem(t.UpdateItemId)

	assert.Error(t.T(), err)
}
