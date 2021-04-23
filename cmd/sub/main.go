package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.sr.ht/~spc/go-log"
	"git.sr.ht/~spc/mqttcli"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/peterbourgon/ff/v3"
)

func main() {
	fs := mqttcli.NewFlagSet("sub", flag.ExitOnError)

	if err := ff.Parse(fs, os.Args[1:], ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser)); err != nil {
		log.Fatalf("failed parse: %v", err)
	}

	if mqttcli.Broker == "" {
		fmt.Println("missing required flag: -broker")
		os.Exit(2)
	}

	if mqttcli.Verbose {
		log.SetLevel(log.LevelInfo)
	}

	if _, ok := os.LookupEnv("MQTTDEBUG"); ok {
		mqtt.DEBUG = log.New(os.Stderr, "[DEBUG] ", log.Flags(), log.CurrentLevel())
	}

	opts := mqttcli.NewClientOptions()

	opts.SetOnConnectHandler(func(c mqtt.Client) {
		for _, topic := range mqttcli.Topics.Values {
			c.Subscribe(topic, byte(mqttcli.QoS), func(c mqtt.Client, m mqtt.Message) {
				log.Printf("[%v] %v", m.Topic(), string(m.Payload()))
			})
			log.Infof("subscribed: %v", topic)
		}
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
