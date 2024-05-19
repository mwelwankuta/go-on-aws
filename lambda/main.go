package main

import (
	"lambdafunc/api/app"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"`
}

func main() {
	myApp := app.NewApp()
	handler := myApp.ApiHandler

	lambda.Start(handler.RegisterUserHandler)
}
