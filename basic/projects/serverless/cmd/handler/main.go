package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/gopro/basic/projects/serverless/internal/handlers"
)

func main() {
	lambda.Start(handlers.HandleRequest)
}
