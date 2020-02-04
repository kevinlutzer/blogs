package main

import "io/ioutil"

const (
	DeviceID   = "test-device"
	Host       = "mqtt.googleapis.com"
	Port       = "8883"
	ProjectID  = "iot-klutzer"
	RegistryID = "devices-klutzer"
	Region     = "us-central1"
	Topic      = "/devices/test-device/config"
	Path       = "certs/"
)

type sslCerts struct {
	RSAPrivate string
	Roots      string
}

func getSSLCerts() (*sslCerts, error) {
	rs, err := ioutil.ReadFile(Path + "roots.pem")
	if err != nil {
		return nil, err
	}

	rscp, err := ioutil.ReadFile(Path + "rsa_private.pem")
	if err != nil {
		return nil, err
	}

	return &sslCerts{
		Roots:      string(rs),
		RSAPrivate: string(rscp),
	}, nil
}
