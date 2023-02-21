package main

import (
	"context"
	"log"
	resources "resources/pkg"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	tag "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

var (
	client *tag.Client
	ses    *sesv2.Client
)

func Handler() {
	results, err := resources.GetTagResources(client, "Owner", "owner")
	if err != nil {
		log.Fatal(err)
	}
	resources.SendSns(ses, results)
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	// NewFromConfig returns a new client from the provided config.
	client = tag.NewFromConfig(cfg)
	ses = sesv2.NewFromConfig(cfg)
}

func main() {
	lambda.Start(Handler)
}
