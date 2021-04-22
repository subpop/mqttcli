package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.sr.ht/~spc/go-log"
	mqttcli "git.sr.ht/~spc/mqttcli"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/peterbourgon/ff/v3"
)

func main() {
	fs := mqttcli.NewFlagSet("sub", flag.ExitOnError)

	if err := ff.Parse(fs, os.Args[1:], ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser)); err != nil {
		log.Fatalf("failed parse: %v", err)
	}

	if mqttcli.Broker == "" || mqttcli.Topic == "" {
		fs.Usage()
		os.Exit(0)
	}

	logLevel, err := log.ParseLevel(mqttcli.LogLevel)
	if err != nil {
		log.Fatalf("cannot parse log level: %v", err)
	}
	log.SetLevel(logLevel)

	if log.CurrentLevel() >= log.LevelDebug {
		mqtt.DEBUG = log.New(os.Stderr, "[DEBUG] ", log.Flags(), log.CurrentLevel())
	}

	opts := mqttcli.NewClientOptions()

	opts.SetOnConnectHandler(func(c mqtt.Client) {
		c.Subscribe(mqttcli.Topic, byte(mqttcli.QoS), func(c mqtt.Client, m mqtt.Message) {
			log.Printf("[%v] %v", m.Topic(), string(m.Payload()))
		})
		log.Infof("subscribed: %v", mqttcli.Topic)
	})

	client := mqtt.NewClient(opts)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGABRT)

	if token := client.Connect(); token.WaitTimeout(30*time.Second) && token.Error() != nil {
		log.Fatalf("connect failed: %v", token.Error())
	}
	log.Infof("connected: %v", mqttcli.Broker)

	<-quit
}
