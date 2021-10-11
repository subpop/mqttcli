package mqttcli

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// NewMQTTClientOptions returns a MQTT client option structure, prepopulated
// with values read from the flag variables. It is an error to call this
// function before parsing flags.
func NewMQTTClientOptions() (*mqtt.ClientOptions, error) {
	opts := mqtt.NewClientOptions()

	opts.SetOrderMatters(false)

	opts.AddBroker(Broker)
	opts.SetClientID(ClientID)
	opts.SetCleanSession(Clean)

	switch {
	case Username != "" && Password != "" && CertFile.Value != "" && KeyFile.Value != "":
		return nil, fmt.Errorf("authentication can only be one of username/password or cert-file/key-file")
	case Username != "" && Password != "":
		opts.SetUsername(Username)
		opts.SetPassword(Password)
	case CertFile.Value != "" && KeyFile.Value != "":
		config := &tls.Config{}

		certData, err := ioutil.ReadFile(CertFile.Value)
		if err != nil {
			return nil, fmt.Errorf("cannot read cert-file: %w", err)
		}

		keyData, err := ioutil.ReadFile(KeyFile.Value)
		if err != nil {
			return nil, fmt.Errorf("cannot read key-file: %w", err)
		}

		cert, err := tls.X509KeyPair(certData, keyData)
		if err != nil {
			return nil, fmt.Errorf("cannot create certificate: %w", err)
		}

		config.Certificates = append(config.Certificates, cert)

		if CARoot != "" {
			pool, err := x509.SystemCertPool()
			if err != nil {
				return nil, fmt.Errorf("cannot get system certificate pool: %w", err)
			}

			data, err := ioutil.ReadFile(CARoot)
			if err != nil {
				return nil, fmt.Errorf("cannot read file: %w", err)
			}

			pool.AppendCertsFromPEM(data)
			config.RootCAs = pool
		}

		opts.SetTLSConfig(config)
	default:
		if Username != "" && Password == "" {
			return nil, fmt.Errorf("password required when using username")
		}
		if Username == "" && Password != "" {
			return nil, fmt.Errorf("username required when using password")
		}
		if CertFile.Value != "" && KeyFile.Value == "" {
			return nil, fmt.Errorf("key-file required when using cert-file")
		}
		if CertFile.Value == "" && KeyFile.Value != "" {
			return nil, fmt.Errorf("cert-file required when using key-file")
		}
	}

	h := make(http.Header)
	for k, v := range Headers.Values {
		h.Add(k, v)
	}
	opts.SetHTTPHeaders(h)

	return opts, nil
}
