package fume

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambdaV2

type Options struct {
	Port int // port to run on (default 8080)
}

func Start(routes *gin.Engine, options *Options) {

	defaults := &Options{
		Port: 8080,
	}

	if options.Port != 0 {
		defaults.Port = options.Port
	}

	if os.Getenv("_HANDLER") != "" {
		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", defaults.Port),
			Handler: routes,
		}
		server.ListenAndServe()
	} else {
		ginLambda = ginadapter.NewV2(routes)
		lambda.Start(Handler)
	}
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}