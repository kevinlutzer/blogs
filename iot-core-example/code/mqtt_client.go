package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"os"
)

const (
	deviceID    = "test-device"
	host        = "mqtt.googleapis.com"
	port        = "8883"
	registryID  = "devices"
	region      = "us-central1"
	configTopic = "/devices/test-device/config"
	certPath    = "certs/"
)

var projectID string

func init() {
	if projectID = os.Getenv("PROJECT_ID"); projectID == "" {
		panic("PROJECT_ID environment variable is required.\n Please start this application by running `PROJECT_ID=<INSERT_PROJECT_ID>`")
	}
}

func getSSLCerts() (rootsCert []byte, clientKey []byte, err error) {
	rootsCert, err = ioutil.ReadFile(certPath + "roots.pem")
	if err != nil {
		return
	}

	clientKey, err = ioutil.ReadFile(certPath + "rsa_private.pem")
	if err != nil {
		return
	}

	return
}

func getTokenString(rsaPrivate string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)

	token.Claims = jwt.StandardClaims{
		Audience:  projectID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(rsaPrivate))
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func getTLSConfig(rootsCert string) *tls.Config {
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM([]byte(rootsCert))

	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{},
		MinVersion:         tls.VersionTLS12,
	}
}

func newClient() (*client, error) {

	clientID := fmt.Sprintf("projects/%v/locations/%v/registries/%v/devices/%v",
		projectID,
		region,
		registryID,
		deviceID,
	)

	roots, clientKey, err := getSSLCerts()
	if err != nil {
		return nil, err
	}

	jwtString, err := getTokenString(string(clientKey))
	if err != nil {
		return nil, err
	}

	tlsConfig := getTLSConfig(string(roots))

	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%v:%v", host, port))
	opts.SetClientID(clientID).SetTLSConfig(tlsConfig)
	opts.SetUsername("unused")
	opts.SetPassword(jwtString)

	mqttClient := MQTT.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &client{
		mqttClient,
	}, nil
}

type client struct {
	mqttClient MQTT.Client
}

func (c *client) Subsribe(topic string, f MQTT.MessageHandler) error {
	if token := c.mqttClient.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *client) Publish(topic string, msg interface{}) error {
	if token := c.mqttClient.Publish(topic, 1, false, msg); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
