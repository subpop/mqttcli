package mqttcli

import (
	"flag"
	"net/http"
)

var (
	Broker   string
	Topic    string
	ClientID string
	Username string
	Password string
	CARoot   string
	QoS      int
	LogLevel string
)
var Headers = http.Header{}

func NewFlagSet(name string, errorHandling flag.ErrorHandling) *flag.FlagSet {
	fs := flag.NewFlagSet(name, errorHandling)

	_ = fs.String("config", "", "")
	fs.StringVar(&Broker, "broker", "", "")
	fs.StringVar(&Topic, "topic", "", "")
	fs.StringVar(&ClientID, "client-id", "", "")
	fs.StringVar(&Username, "username", "", "")
	fs.StringVar(&Password, "password", "", "")
	fs.StringVar(&CARoot, "ca-root", "", "")
	fs.IntVar(&QoS, "qos", 0, "")
	fs.StringVar(&LogLevel, "log-level", "error", "")
	fs.Var(&HeadersValue{Headers: Headers}, "headers", "")

	return fs
}
