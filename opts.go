package mqttcli

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewClientOptions() *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(Broker)
	opts.SetClientID(ClientID)
	opts.SetUsername(Username)
	opts.SetPassword(Password)
	opts.SetHTTPHeaders(Headers)

	if CARoot != "" {
		tlsConfig := &tls.Config{}
		pool, err := x509.SystemCertPool()
		if err != nil {
			log.Fatalf("cannot get system certificate pool: %v", err)
		}

		data, err := ioutil.ReadFile(CARoot)
		if err != nil {
			log.Fatalf("cannot read file: %v", err)
		}
		pool.AppendCertsFromPEM(data)
		tlsConfig.RootCAs = pool
		opts.SetTLSConfig(tlsConfig)
	}
	return opts
}
