package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

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

func newClient(certs *sslCerts) (*client, error) {

	clientID := fmt.Sprintf("projects/%v/locations/%v/registries/%v/devices/%v",
		projectID,
		region,
		registryID,
		deviceID,
	)

	jwtString, err := getTokenString(certs.RSAPrivate)
	if err != nil {
		return nil, err
	}

	tlsConfig := getTLSConfig(certs.Roots)

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

func (c *client) Publish(topic, msg string) error {
	if token := c.mqttClient.Publish(topic, 2, false, msg); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
