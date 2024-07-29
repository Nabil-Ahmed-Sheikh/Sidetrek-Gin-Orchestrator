package awsconfig

import (
	"GinProject/app/config/env"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func MustGetCertAndKeyFromSecret(cfg env.Config) (cert string, key string) {
	var err error

	cert, err = getSecret(cfg, cfg.Certs.Cert)
	if err != nil {
		log.Fatalln("unable to get cert", err)
	}

	key, err = getSecret(cfg, cfg.Certs.Key)
	if err != nil {
		log.Fatalln("unable to get key", err)
	}

	return cert, key
}

func getSecret(cfg env.Config, secretName string) (string, error) {
	var secretString string

	awsCfg := LoadConfig(cfg)

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(awsCfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		return secretString, err
	}

	// Decrypts secret using the associated KMS key.
	secretString = *result.SecretString

	return secretString, nil
}
