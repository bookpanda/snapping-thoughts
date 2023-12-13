package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

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

	// item := Item{
	// 	Id:      "1",
	// 	IsUsed:  "no",
	// 	Message: "ok bro",
	// }

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
