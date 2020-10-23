package graphql

// RequestInput represents the graphql request input.
type RequestInput struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
}
