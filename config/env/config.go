package env

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	TEMPORAL_HOSTPORT  = "TEMPORAL_HOSTPORT"
	TEMPORAL_NAMESPACE = "TEMPORAL_NAMESPACE"
	TEMPORAL_TASKQUEUE = "TEMPORAL_TASKQUEUE"

	TERRAFORM_STATE_DIR = "TERRAFORM_STATE_DIR"

	TERRAFORM_STATE_BUCKET = "TERRAFORM_STATE_BUCKET"
	TERRAFORM_STATE_REGION = "TERRAFORM_STATE_REGION"
	TERRAFORM_STATE_DYNAMO = "TERRAFORM_STATE_DYNAMO"

	COGNITO_ISS       = "COGNITO_ISS"
	COGNITO_USER_POOL = "COGNITO_USER_POOL"
	COGNITO_REGION    = "COGNITO_REGION"

	TEMPORAL_CA_CERT = "TEMPORAL_CA_CERT"
	TEMPORAL_CA_KEY  = "TEMPORAL_CA_KEY"

	GITHUB_TOKEN = "GITHUB_TOKEN"
)

type (
	Config struct {
		TemporalHostPort  string
		TemporalNamespace string
		TemporalTaskQueue string
		TfState           TfState
		Cognito           Cognito
		GithubToken       string
		Certs             Certs
	}

	Cognito struct {
		Iss      string
		UserPool string
		Region   string
	}

	TfState struct {
		Region   string // aws region
		Bucket   string // s3 bucket name
		DynamoDB string
	}

	Certs struct {
		Cert string
		Key  string
	}
)

func MustGetConfig() Config {
	godotenv.Load()
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	return Config{
		TemporalHostPort:  getOrDefault(TEMPORAL_HOSTPORT, "127.0.0.1:7233"),
		TemporalNamespace: getOrDefault(TEMPORAL_NAMESPACE, "default"),
		TemporalTaskQueue: getOrDefault(TEMPORAL_TASKQUEUE, "temporal-terraform-demo"),
		TfState:           getTfState(),
		Cognito:           getCognito(),
		GithubToken:       getOrDefault(GITHUB_TOKEN, ""),
		Certs:             getCerts(),
	}
}

func getOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getTfState() TfState {
	return TfState{
		Region:   getOrDefault(TERRAFORM_STATE_REGION, "us-west-1"),
		Bucket:   getOrDefault(TERRAFORM_STATE_BUCKET, "sidetrek-infra-tfstate"),
		DynamoDB: getOrDefault(TERRAFORM_STATE_DYNAMO, "sidetrek-tf-lockid"),
	}
}

func getCognito() Cognito {
	return Cognito{
		Iss:      getOrDefault(COGNITO_ISS, ""),
		UserPool: getOrDefault(COGNITO_USER_POOL, ""),
		Region:   getOrDefault(COGNITO_REGION, ""),
	}
}

func getCerts() Certs {
	return Certs{
		Cert: getOrDefault(TEMPORAL_CA_CERT, "sidetrek-tcld-ca.crt"),
		Key:  getOrDefault(TEMPORAL_CA_KEY, "sidetrek-tcld-ca.key"),
	}
}
