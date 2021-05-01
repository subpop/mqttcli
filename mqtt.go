package mqttcli

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// NewMQTTClientOptions returns a MQTT client option structure, prepopulated
// with values read from the flag variables. It is an error to call this
// function before parsing flags.
func NewMQTTClientOptions() (*mqtt.ClientOptions, error) {
	opts := mqtt.NewClientOptions()

	opts.AddBroker(Broker)
	opts.SetClientID(ClientID)
	opts.SetUsername(Username)
	opts.SetPassword(Password)

	h := make(http.Header)
	for k, v := range Headers.Values {
		h.Add(k, v)
	}
	opts.SetHTTPHeaders(h)

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

	return opts, nil
}
