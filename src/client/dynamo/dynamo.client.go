package dynamo

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/bookpanda/snapping-thoughts/src/model/dynamo"

	"github.com/rs/zerolog/log"
)

type DynamoDBClient struct {
	db        Database
	tableName string
}

type Database interface {
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
}

func NewDynamoDBClient(db Database, tableName string) *DynamoDBClient {
	return &DynamoDBClient{
		db,
		tableName,
	}
}

func (c *DynamoDBClient) CreateItem(item *dynamo.Item) error {
	log.Info().Str("twitterClient", "CreateItem")
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatal().Str("twitterClient", "Got error marshalling new item").Err(err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(c.tableName),
	}

	_, err = c.db.PutItem(input)
	if err != nil {
		log.Fatal().Str("twitterClient", "Got error calling PutItem").Err(err)
	}
	log.Info().Msgf("Successfully added item with id " + item.Id + " to table " + c.tableName)

	return nil
}

func (c *DynamoDBClient) GetItem() (*dynamo.Item, error) {
	log.Info().Str("twitterClient", "GetItem")
	proj := expression.NamesList(expression.Name("Message"), expression.Name("Id"))
	filt := expression.Name("IsUsed").Equal(expression.Value("no"))
	expr, err := expression.NewBuilder().WithProjection(proj).WithFilter(filt).Build()
	if err != nil {
		log.Fatal().Str("twitterClient", "Got error building expression").Err(err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(c.tableName),
		Limit:                     aws.Int64(1),
	}

	result, err := c.db.Scan(params)
	if err != nil {
		log.Fatal().Str("twitterClient", "Got error calling Scan").Err(err)
	}

	if len(result.Items) == 0 {
		log.Info().Str("twitterClient", "Could not find unused item")
		return nil, nil
	}

	item := dynamo.Item{}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	log.Info().Msgf("Successfully scanned item with id " + item.Id + " from table " + c.tableName)

	return &item, nil
}

func (c *DynamoDBClient) GetItemWithId(id string) (*dynamo.Item, error) {
	log.Info().Str("twitterClient", "GetItemWithId").Str("id: ", id)

	params := &dynamodb.GetItemInput{
		TableName: aws.String(c.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
	}

	result, err := c.db.GetItem(params)
	if err != nil {
		log.Fatal().Str("twitterClient", "Got error calling GetItem").Err(err)
	}

	if result == nil {
		log.Info().Str("twitterClient", "Could not find item with id: "+id)
		return nil, nil
	}

	item := dynamo.Item{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	log.Info().Msgf("Successfully got item with id " + item.Id + " from table " + c.tableName)

	return &item, nil
}

func (c *DynamoDBClient) UpdateItem(id string) error {
	log.Info().Str("twitterClient", "Updating item with id: "+id)
	key := map[string]*dynamodb.AttributeValue{
		"Id": {
			S: aws.String(id),
		},
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":u": {
				S: aws.String(time.Now().String()),
			},
		},
		TableName:        aws.String(c.tableName),
		Key:              key,
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set IsUsed = :u"),
	}

	_, err := c.db.UpdateItem(input)
	if err != nil {
		log.Fatal().Str("twitterClient", "Got error calling UpdateItem").Err(err)
	}
	log.Info().Msgf("Successfully updated item with id " + id + " to table " + c.tableName)

	return nil
}
