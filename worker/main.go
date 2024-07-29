package main

import (
	"GinProject/app/config/env"
	"GinProject/app/worker/workflows"
	"crypto/tls"
	"encoding/base64"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	// "github.com/joho/godotenv"
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
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	temporalWorker := worker.New(c, cfg.TemporalTaskQueue, worker.Options{
		WorkerStopTimeout: 30 * time.Second,
	})

	log.Print("registering workflows")
	workflows.Register(temporalWorker)

	if err := temporalWorker.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("hostPort.", hostPort)
		log.Fatalln("unable to start Worker", err)
	}
}
