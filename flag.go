package mqttcli

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/sgreben/flagvar"
)

var (
	Broker   string
	Topics   flagvar.Strings
	ClientID string
	Username string
	Password string
	CARoot   string
	QoS      int
	Verbose  bool
	Headers  flagvar.AssignmentsMap
)

// GlobalFlagSet returns a new flag set configured with flags common to publish
// and subscribe clients.
func GlobalFlagSet(name string, errorHandling flag.ErrorHandling) *flag.FlagSet {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	if name == "" {
		name = filepath.Base(os.Args[0])
	}

	fs := flag.NewFlagSet(name, errorHandling)

	fs.StringVar(&Broker, "broker", "", "broker address (should be in the form of a `URL`)")
	fs.Var(&Topics, "topic", "topic to publish or subscribe to\n(can be specified multiple times)")
	fs.StringVar(&ClientID, "client-id", hostname+"-"+randomString(6), "unique identifier for this client")
	fs.StringVar(&Username, "username", "", "authenticate with a username")
	fs.StringVar(&Password, "password", "", "authenticate with a password")
	fs.StringVar(&CARoot, "ca-root", "", "path to a `file` containing CA certificates")
	fs.IntVar(&QoS, "qos", 0, "quality of service for messages")
	fs.BoolVar(&Verbose, "verbose", false, "increase output")
	fs.Var(&Headers, "header", "set an HTTP header (in `KEY=VALUE` form)\n(can be specified multiple times)")

	return fs
}
