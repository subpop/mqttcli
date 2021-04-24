package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/sgreben/flagvar"
)

func main() {

	var (
		rootFs = flag.NewFlagSet("mqttcli", flag.ExitOnError)
		pubFs  = flag.NewFlagSet("mqtt pub", flag.ExitOnError)
		subFs  = flag.NewFlagSet("mqtt sub", flag.ExitOnError)
		opts   = mqtt.NewClientOptions()

		_        = rootFs.String("config", "", "path to `file` containing configuration values (optional)")
		broker   string
		topics   flagvar.Strings
		clientID string
		username string
		password string
		caRoot   string
		qos      int
		verbose  bool
		headers  flagvar.AssignmentsMap

		message  string
		retained bool
	)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	rootFs.StringVar(&broker, "broker", "", "broker address (should be in the form of a `URL`)")
	rootFs.Var(&topics, "topic", "topic to publish or subscribe to\n(can be specified multiple times)")
	rootFs.StringVar(&clientID, "client-id", hostname+"-"+randomString(6), "unique identifier for this client")
	rootFs.StringVar(&username, "username", "", "authenticate with a username")
	rootFs.StringVar(&password, "password", "", "authenticate with a password")
	rootFs.StringVar(&caRoot, "ca-root", "", "path to a `file` containing CA certificates")
	rootFs.IntVar(&qos, "qos", 0, "quality of service for messages")
	rootFs.BoolVar(&verbose, "verbose", false, "increase output")
	rootFs.Var(&headers, "header", "set an HTTP header (in `KEY=VALUE` form)\n(can be specified multiple times)")

	pubFs.StringVar(&message, "message", "", "message payload")
	pubFs.BoolVar(&retained, "retained", false, "retain message on the broker")

	pub := &ffcli.Command{
		Name:       "pub",
		ShortUsage: "mqttcli [global flags] pub [flags]",
		ShortHelp:  "Publish a message to the broker.",
		FlagSet:    pubFs,
		Exec: func(ctx context.Context, args []string) error {
			client := mqtt.NewClient(opts)
			if token := client.Connect(); token.WaitTimeout(30*time.Second) && token.Error() != nil {
				log.Fatalf("connect failed: %v", token.Error())
			}
			log.Printf("connected: %v", broker)

			for _, topic := range topics.Values {
				token := client.Publish(topic, byte(qos), retained, message)
				if token.Wait() && token.Error() != nil {
					log.Fatalf("publish failed: %v", token.Error())
				}
				log.Printf("published: [%v] %v", topic, message)
			}

			return nil
		},
	}
	sub := &ffcli.Command{
		Name:       "sub",
		ShortUsage: "mqttcli [global flags] sub",
		ShortHelp:  "Subscribe to topics on the broker.",
		FlagSet:    subFs,
		Exec: func(ctx context.Context, args []string) error {
			opts.SetOnConnectHandler(func(c mqtt.Client) {
				for _, topic := range topics.Values {
					c.Subscribe(topic, byte(qos), func(c mqtt.Client, m mqtt.Message) {
						log.Printf("[%v] %v", m.Topic(), string(m.Payload()))
					})
					log.Printf("subscribed: %v", topic)
				}
			})

			client := mqtt.NewClient(opts)
			if token := client.Connect(); token.WaitTimeout(30*time.Second) && token.Error() != nil {
				log.Fatalf("connect failed: %v", token.Error())
			}
			log.Printf("connected: %v", broker)

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGTERM, syscall.SIGABRT)

			<-quit

			return nil
		},
	}

	root := &ffcli.Command{
		ShortUsage:  "mqttcli [global flags] <subcommand>",
		FlagSet:     rootFs,
		Options:     []ff.Option{ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser)},
		Subcommands: []*ffcli.Command{pub, sub},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}

	if err := root.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)

	h := make(http.Header)
	for k, v := range headers.Values {
		h.Add(k, v)
	}
	opts.SetHTTPHeaders(h)

	if caRoot != "" {
		tlsConfig := &tls.Config{}
		pool, err := x509.SystemCertPool()
		if err != nil {
			log.Fatalf("cannot get system certificate pool: %v", err)
		}

		data, err := ioutil.ReadFile(caRoot)
		if err != nil {
			log.Fatalf("cannot read file: %v", err)
		}
		pool.AppendCertsFromPEM(data)
		tlsConfig.RootCAs = pool
		opts.SetTLSConfig(tlsConfig)
	}

	if err := root.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
