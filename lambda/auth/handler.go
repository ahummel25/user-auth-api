package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"

	g "github.com/src/user-auth-api/graphql"
	"github.com/src/user-auth-api/utils"
)

// LambdaHandler is our lambda handler invoked by the `lambda.Start` function call
func LambdaHandler(
	request events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	var (
		buf      *bytes.Buffer
		err      error
		gqlErr   gqlerrors.FormattedError
		payload  g.RequestInput
		response = events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			IsBase64Encoded: false,
		}
		responseBody []byte
		result       *graphql.Result
	)

	if err = json.Unmarshal([]byte(request.Body), &payload); err != nil {
		response = utils.BuildErrorResponse(response, err.Error())

		return response, nil
	}

	if result, gqlErr = g.ExecuteQuery(payload); gqlErr.Message != "" {
		response = utils.BuildErrorResponse(response, gqlErr.Message)

		return response, nil
	}

	responseBody, _ = json.Marshal(map[string]interface{}{
		"data": result.Data,
	})

	buf = bytes.NewBuffer(responseBody)

	response.StatusCode = http.StatusOK
	response.Body = buf.String()

	return response, nil
}
