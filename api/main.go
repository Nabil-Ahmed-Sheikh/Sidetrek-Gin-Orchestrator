package main

import (
	"crypto/tls"
	"encoding/base64"
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.uber.org/zap"

	// "go.temporal.io/sdk/worker"

	"GinProject/app/api/handlers"
	"GinProject/app/config/env"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	// "github.com/joho/godotenv"
)

var (
	binding string // port
	isDev   bool   // checks if is dev
)

func main() {
	cfg := env.MustGetConfig()

	godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	clientCert := os.Getenv("TEMPORAL_CA_CERT")
	clientKey := os.Getenv("TEMPORAL_CA_KEY")
	hostPort := os.Getenv("TEMPORAL_HOSTPORT")
	namespace := os.Getenv("TEMPORAL_NAMESPACE")

	var cert tls.Certificate
	var err error

	env := os.Getenv("ENVIRONMENT")
	if env == "local" {
		cert, err = tls.LoadX509KeyPair(clientCert, clientKey)
		if err != nil {
			log.Fatalln("Unable to load cert and key pair.", err)
		}
	} else {
		// Decode the base64 strings
		clientCertPEM, err := base64.StdEncoding.DecodeString(clientCert)
		if err != nil {
			log.Fatalln("Unable to decode client certificate base64 string.", err)
		}
		clientKeyPEM, err := base64.StdEncoding.DecodeString(clientKey)
		if err != nil {
			log.Fatalln("Unable to decode client key base64 string.", err)
		}

		cert, err = tls.X509KeyPair(clientCertPEM, clientKeyPEM)
		if err != nil {
			log.Fatalln("Client Cert.", clientCert)
			log.Fatalln("Client Key.", clientKey)
			log.Fatalln("Unable to load cert and key pair.", err)
		}
	}

	c, err := client.NewLazyClient(client.Options{
		HostPort:  hostPort,
		Namespace: namespace,
		ConnectionOptions: client.ConnectionOptions{
			TLS: &tls.Config{Certificates: []tls.Certificate{cert}},
		},
	})

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	api := r.Group("/api")

	handlers.NewWorkflowHandler(api, c, cfg.TemporalTaskQueue, logger.Sugar())

	if err := r.Run(":3000"); err != nil {
		panic(err)
	}

	// server.Run(":3000") // listen and serve on 0.0.0.0:8080
}
