package main

import (
	"flag"
	"os"
	"time"

	"git.sr.ht/~spc/go-log"
	"git.sr.ht/~spc/mqttcli"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/peterbourgon/ff/v3"
)

func main() {
	fs := mqttcli.NewFlagSet("sub", flag.ExitOnError)

	var message = fs.String("message", "", "")
	var retained = fs.Bool("retained", false, "")

	if err := ff.Parse(fs, os.Args[1:], ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser)); err != nil {
		log.Fatalf("failed parse: %v", err)
	}

	logLevel, err := log.ParseLevel(mqttcli.LogLevel)
	if err != nil {
		log.Fatalf("cannot parse log level: %v", err)
	}
	log.SetLevel(logLevel)

	if log.CurrentLevel() >= log.LevelDebug {
		mqtt.DEBUG = log.New(os.Stderr, "[DEBUG] ", log.Flags(), log.CurrentLevel())
	}

	if mqttcli.Broker == "" || mqttcli.Topic == "" || *message == "" {
		fs.Usage()
		os.Exit(0)
	}

	opts := mqttcli.NewClientOptions()

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.WaitTimeout(30*time.Second) && token.Error() != nil {
		log.Fatalf("connect failed: %v", token.Error())
	}
	log.Infof("connected: %v", mqttcli.Broker)

	token := client.Publish(mqttcli.Topic, byte(mqttcli.QoS), *retained, *message)
	if token.Wait() && token.Error() != nil {
		log.Fatalf("publish failed: %v", token.Error())
	}
	log.Infof("published: [%v] %v", mqttcli.Topic, *message)
}
