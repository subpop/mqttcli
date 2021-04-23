package mqttcli

import (
	"flag"

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

func NewFlagSet(name string, errorHandling flag.ErrorHandling) *flag.FlagSet {
	fs := flag.NewFlagSet(name, errorHandling)

	_ = fs.String("config", "", "")
	fs.StringVar(&Broker, "broker", "", "")
	fs.Var(&Topics, "topic", "")
	fs.StringVar(&ClientID, "client-id", "", "")
	fs.StringVar(&Username, "username", "", "")
	fs.StringVar(&Password, "password", "", "")
	fs.StringVar(&CARoot, "ca-root", "", "")
	fs.IntVar(&QoS, "qos", 0, "")
	fs.BoolVar(&Verbose, "verbose", false, "")
	fs.Var(&Headers, "header", "")

	return fs
}
