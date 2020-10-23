package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/ahummel25/user-auth-api/lambda/graphql"
)

func main() {
	lambda.Start(graphql.LambdaHandler)
}
