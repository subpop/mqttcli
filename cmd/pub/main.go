package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"git.sr.ht/~spc/go-log"
	"git.sr.ht/~spc/mqttcli"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/peterbourgon/ff/v3"
)

func main() {
	fs := mqttcli.NewFlagSet("pub", flag.ExitOnError)

	var message = fs.String("message", "", "message payload")
	var retained = fs.Bool("retained", false, "retain message on the broker")

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

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.WaitTimeout(30*time.Second) && token.Error() != nil {
		log.Fatalf("connect failed: %v", token.Error())
	}
	log.Infof("connected: %v", mqttcli.Broker)

	for _, topic := range mqttcli.Topics.Values {
		token := client.Publish(topic, byte(mqttcli.QoS), *retained, *message)
		if token.Wait() && token.Error() != nil {
			log.Fatalf("publish failed: %v", token.Error())
		}
		log.Infof("published: [%v] %v", topic, *message)
	}
}
