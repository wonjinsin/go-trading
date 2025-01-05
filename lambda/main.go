package main

import (
	"context"
	"fmt"
	"magmar/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
)

func main() {
	fmt.Println("Hello, World!")
	lambda.Start(Handler)
}

// Handler ...
func Handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	apiGatewayRequest := events.APIGatewayProxyRequest{
		Path:       req.RawPath,
		HTTPMethod: req.RequestContext.HTTP.Method,
		Headers:    req.Headers,
		Body:       req.Body,
	}

	if apiGatewayRequest.Path == "" {
		apiGatewayRequest.Path = "/"
	}

	if !strings.HasPrefix(apiGatewayRequest.Path, "/") {
		apiGatewayRequest.Path = "/" + apiGatewayRequest.Path
	}

	e := http.EchoHandler()
	echoLambda := echoadapter.New(e)
	response, err := echoLambda.ProxyWithContext(ctx, apiGatewayRequest)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error: %v", err),
		}, nil
	}

	// API Gateway 응답을 Function URL 응답으로 변환
	return events.LambdaFunctionURLResponse{
		StatusCode: response.StatusCode,
		Headers:    response.Headers,
		Body:       response.Body,
	}, nil
}
