package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

func dataFromDynamoDB() ([]Book, error) {
	fmt.Printf("DynamoDBから全ての書籍データを取得します。")
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")},
	)
	fmt.Println("DynamoDBとの接続に成功しました。")

	svc := dynamodb.New(sess)

	tableName := "books-table"

	fmt.Println("DynamoDBから全ての書籍データを取得します。")
	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println("DynamoDBから全ての書籍データを取得しました。")

	books := []Book{}

	for _, i := range result.Items {
		book := Book{}
		err = dynamodbattribute.UnmarshalMap(i, &book)

		if err != nil {
			panic(fmt.Sprintf("レコードのアンマーシャルに失敗しました, %v", err))
		}

		books = append(books, book)
	}

	fmt.Println("DynamoDBから全ての書籍データが正常に取得できました。")
	return books, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}

	// DynamoDBと接続するときはコメントアウトを外す
	books, _ := dataFromDynamoDB()

	// // DynamoDBのテーブルを作成するまでのダミーデータ
	// books := []Book{
	// 	Book{ID: "1", Title: "Book1", Author: "Author1"},
	// 	Book{ID: "2", Title: "Book2", Author: "Author2"},
	// 	Book{ID: "3", Title: "Book3", Author: "Author3"},
	// }

	bytes, err := json.Marshal(books)
	if err != nil {
		fmt.Println(err.Error())
	}

	return events.APIGatewayProxyResponse{
		Body:       string(bytes),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
