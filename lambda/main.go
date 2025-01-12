package main

import (
	"context"
	"fmt"
	"log"
	"magmar/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
)

var (
	echoLambda  *echoadapter.EchoLambda
	initialized bool
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("Starting initialization...")
}

func main() {
	e := http.EchoHandler()
	if !initialized {
		echoLambda = echoadapter.New(e)
		initialized = true
		log.Println("Echo initialized, starting handler...")
	}
	log.Println("Starting Lambda handler...")
	lambda.Start(Handler)
}

// Handler ...
func Handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	log.Printf("Handling request - path: %s, method: %s, requestID: %s\n",
		req.RawPath,
		req.RequestContext.HTTP.Method,
		req.RequestContext.RequestID)

	select {
	case <-ctx.Done():
		log.Printf("Context cancelled: %v\n", ctx.Err())
		return events.LambdaFunctionURLResponse{
			StatusCode: 504,
			Body:       "Request timeout",
		}, ctx.Err()
	default:
	}

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
