package client

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"log"
)

type DynamoDBClient struct {
	client    *dynamodb.DynamoDB
	tableName string
}

func NewDynamoDBClient(tableName string) *DynamoDBClient {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := dynamodb.New(sess)

	return &DynamoDBClient{
		client,
		tableName,
	}
}

type Item struct {
	Id      string
	IsUsed  string //no or date
	Message string
}

func (c *DynamoDBClient) CreateItem(item Item) error {
	av, err := dynamodbattribute.MarshalMap(item)
	log.Println(av)
	if err != nil {
		log.Fatalf("Got error marshalling new item: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(c.tableName),
	}

	_, err = c.client.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
	log.Println("Successfully added '" + item.Message + " to table " + c.tableName)

	return nil
}

func (c *DynamoDBClient) GetItem() (*string, error) {
	proj := expression.NamesList(expression.Name("Message"))
	filt := expression.Name("IsUsed").Equal(expression.Value("no"))
	expr, err := expression.NewBuilder().WithProjection(proj).WithFilter(filt).Build()
	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(c.tableName),
		Limit:                     aws.Int64(1),
	}

	result, err := c.client.Scan(params)
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	if result.Items[0] == nil {
		log.Println("Could not find unused item")
		return nil, nil
	}

	item := Item{}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return &item.Message, nil
}
