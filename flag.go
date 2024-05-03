package mqttcli

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/sgreben/flagvar"
)

var (
	Broker               string
	Topics               flagvar.Strings
	ClientID             string
	Username             string
	Password             string
	CARoot               string
	QoS                  int
	Verbose              bool
	Clean                bool
	Headers              flagvar.AssignmentsMap
	CertFile             flagvar.File
	KeyFile              flagvar.File
	ConnectRetry         bool
	ConnectRetryInterval time.Duration
	AutoReconnect        bool
	PrintVersion         bool
	TLSALPN              flagvar.Strings
)

const Version = "0.2.6"

// GlobalFlagSet returns a new flag set configured with flags common to publish
// and subscribe clients.
func GlobalFlagSet(name string, errorHandling flag.ErrorHandling) *flag.FlagSet {
	if name == "" {
		name = filepath.Base(os.Args[0])
	}

	fs := flag.NewFlagSet(name, errorHandling)

	fs.StringVar(&Broker, "broker", "", "broker address (should be in the form of a `URL`)")
	fs.Var(&Topics, "topic", "topic to publish or subscribe to\n(can be specified multiple times)")
	fs.StringVar(&ClientID, "client-id", randomString(23), "unique identifier for this client")
	fs.StringVar(&Username, "username", "", "authenticate with a username")
	fs.StringVar(&Password, "password", "", "authenticate with a password")
	fs.StringVar(&CARoot, "ca-root", "", "path to a `file` containing CA certificates")
	fs.IntVar(&QoS, "qos", 0, "quality of service for messages")
	fs.BoolVar(&Verbose, "verbose", false, "increase output")
	fs.BoolVar(&Clean, "clean", true, "discard any pending messages from the broker")
	fs.Var(&Headers, "header", "set an HTTP header (in `KEY=VALUE` form)\n(can be specified multiple times)")
	fs.Var(&CertFile, "cert-file", "authenticate with a certificate")
	fs.Var(&KeyFile, "key-file", "authenticate with a private key")
	fs.BoolVar(&ConnectRetry, "connect-retry", false, "automatically retry initial connection to broker")
	fs.DurationVar(&ConnectRetryInterval, "connect-retry-interval", 30*time.Second, "wait `DURATION` between initial connection attempts")
	fs.BoolVar(&AutoReconnect, "auto-reconnect", false, "automatically reconnect when connection is lost")
	fs.BoolVar(&PrintVersion, "version", false, "print version")
	fs.Var(&TLSALPN, "tls-alpn", "ALPN value to include in the TLS handshake\n(can be specified multiple times)")

	return fs
}
