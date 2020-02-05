package main

import (
	"io/ioutil"
	"os"
)

const (
	deviceID       = "test-device"
	host           = "mqtt.googleapis.com"
	port           = "8883"
	registryID     = "devices"
	region         = "us-central1"
	configTopic    = "/devices/test-device/config"
	telemetryTopic = "/devices/test-device/events"
	certPath       = "certs/"
)

var projectID string

func init() {
	if projectID = os.Getenv("PROJECT_ID"); projectID == "" {
		panic("PROJECT_ID environment variable is required.\n Please start this application by running `PROJECT_ID=<INSERT_PROJECT_ID>`")
	}
}

type sslCerts struct {
	RSAPrivate string
	Roots      string
}

func getSSLCerts() (*sslCerts, error) {
	rs, err := ioutil.ReadFile(certPath + "roots.pem")
	if err != nil {
		return nil, err
	}

	rscp, err := ioutil.ReadFile(certPath + "rsa_private.pem")
	if err != nil {
		return nil, err
	}

	return &sslCerts{
		Roots:      string(rs),
		RSAPrivate: string(rscp),
	}, nil
}
