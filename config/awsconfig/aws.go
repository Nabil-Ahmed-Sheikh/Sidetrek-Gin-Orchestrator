package awsconfig

import (
	"GinProject/app/config/env"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func LoadConfig(cfg env.Config) aws.Config {
	awsConfig, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(cfg.TfState.Region))
	if err != nil {
		log.Fatal("unable to load aws config")
	}
	return awsConfig
}
