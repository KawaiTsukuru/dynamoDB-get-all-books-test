package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// DynamoDBクライアントの作成
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// データを追加するテーブル名
	tableName := "kawai-dynamoDB-books-test"

	// リクエストボディからBookオブジェクトをパース
	var book Book
	err := json.Unmarshal([]byte(request.Body), &book)
	if err != nil {
		log.Printf("リクエストボディのパースエラー: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	// データをDynamoDBに追加
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id":     {S: aws.String(book.ID)},
			"title":  {S: aws.String(book.Title)},
			"author": {S: aws.String(book.Author)},
		},
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Printf("データの追加エラー: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	log.Println("データをDynamoDBに追加しました")

	// レスポンスを返す
	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "データをDynamoDBに追加しました",
	}

	return response, nil
}

func main() {
	lambda.Start(handler)
}
