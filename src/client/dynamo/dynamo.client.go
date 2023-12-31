package dynamo

import (
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
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatal().Str("dynamoClient", "Got error marshalling new item").Err(err)
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(c.tableName),
	}

	_, err = c.db.PutItem(input)
	if err != nil {
		log.Fatal().Str("dynamoClient", "Got error calling PutItem").Err(err)
		return err
	}
	log.Info().Msgf("Successfully added item with id " + item.Id + " to table " + c.tableName)

	return nil
}

func (c *DynamoDBClient) GetItem() (*dynamo.Item, error) {
	proj := expression.NamesList(expression.Name("Message"), expression.Name("Id"))
	filt := expression.Name("IsUsed").Equal(expression.Value("no"))
	expr, err := expression.NewBuilder().WithProjection(proj).WithFilter(filt).Build()
	if err != nil {
		log.Fatal().Str("dynamoClient", "Got error building expression").Err(err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(c.tableName),
	}

	result, err := c.db.Scan(params)
	if err != nil {
		log.Fatal().Str("dynamoClient", "Got error calling Scan").Err(err)
		return nil, err
	}

	if len(result.Items) == 0 {
		log.Info().Str("dynamoClient", "GetItem").Msg("Could not find unused item")
		return nil, nil
	}

	item := dynamo.Item{}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &item)
	if err != nil {
		log.Fatal().Str("dynamoClient", "Got error unmarshalling item").Err(err)
		return nil, err
	}
	log.Info().Msgf("Successfully scanned item with id " + item.Id + " from table " + c.tableName)

	return &item, nil
}

func (c *DynamoDBClient) GetItemWithId(id string) (*dynamo.Item, error) {
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
		log.Fatal().Str("dynamoClient", "Got error calling GetItem").Err(err)
		return nil, err
	}

	if result == nil {
		log.Info().Str("dynamoClient", "GetItemWithId").Msgf("Could not find item with id: " + id)
		return nil, nil
	}

	item := dynamo.Item{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		log.Fatal().Str("dynamoClient", "Got error unmarshalling item").Err(err)
		return nil, err
	}
	log.Info().Msgf("Successfully got item with id " + item.Id + " from table " + c.tableName)

	return &item, nil
}

func (c *DynamoDBClient) UpdateItem(time time.Time, id string) error {
	key := map[string]*dynamodb.AttributeValue{
		"Id": {
			S: aws.String(id),
		},
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":u": {
				S: aws.String(time.String()),
			},
		},
		TableName:        aws.String(c.tableName),
		Key:              key,
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set IsUsed = :u"),
	}

	_, err := c.db.UpdateItem(input)
	if err != nil {
		log.Fatal().Str("dynamoClient", "Got error calling UpdateItem").Err(err)
		return err
	}
	log.Info().Msgf("Successfully updated item with id " + id + " to table " + c.tableName)

	return nil
}
