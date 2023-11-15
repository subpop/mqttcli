package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/peterbourgon/ff/v3"
	"github.com/subpop/mqttcli"
)

var retained bool

func main() {
	fs := mqttcli.GlobalFlagSet("pub", flag.ExitOnError)

	_ = fs.String("config", "", "path to `file` containing configuration values (optional)")
	fs.BoolVar(&retained, "retained", false, "retain message on the broker")

	if err := ff.Parse(fs, os.Args[1:], ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser)); err != nil {
		log.Fatalf("failed to parse flags: %v", err)
	}

	if mqttcli.PrintVersion {
		fmt.Println(mqttcli.Version)
		os.Exit(0)
	}

	message, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read message: %v", err)
	}

	opts, err := mqttcli.NewMQTTClientOptions()
	if err != nil {
		log.Fatalf("failed to configure MQTT: %v", err)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.WaitTimeout(30*time.Second) && token.Error() != nil {
		log.Fatalf("connect failed: %v", token.Error())
	}
	log.Printf("connected: %v", mqttcli.Broker)

	for _, topic := range mqttcli.Topics.Values {
		token := client.Publish(topic, byte(mqttcli.QoS), retained, message)
		if token.Wait() && token.Error() != nil {
			log.Fatalf("publish failed: %v", token.Error())
		}
		log.Printf("published: [%v] %v", topic, message)
	}
}
