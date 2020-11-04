package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

// BuildErrorResponse will build a common error response.
func BuildErrorResponse(response events.APIGatewayProxyResponse, errorMsg string) events.APIGatewayProxyResponse {
	statusCode := http.StatusBadRequest

	if errorMsg == "Access denied!" {
		statusCode = http.StatusForbidden
	}

	errorBody, _ := json.Marshal(map[string]interface{}{
		"message": errorMsg,
	})

	errBuf := bytes.NewBuffer(errorBody)
	response.Body = errBuf.String()
	response.StatusCode = statusCode

	return response
}

// FileExists will verify if the file exists in the specified path.
func FileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
