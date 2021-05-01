package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.sr.ht/~spc/mqttcli"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/peterbourgon/ff/v3"
)

func main() {
	fs := mqttcli.GlobalFlagSet("", flag.ExitOnError)

	_ = fs.String("config", "", "path to `file` containing configuration values (optional)")

	ff.Parse(fs, os.Args[1:], ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))

	opts, err := mqttcli.NewMQTTClientOptions()
	if err != nil {
		log.Fatalf("failed to configure MQTT: %v", err)
	}

	opts.SetOnConnectHandler(func(c mqtt.Client) {
		for _, topic := range mqttcli.Topics.Values {
			c.Subscribe(topic, byte(mqttcli.QoS), func(c mqtt.Client, m mqtt.Message) {
				log.Printf("[%v] %v", m.Topic(), string(m.Payload()))
			})
			log.Printf("subscribed: %v", topic)
		}
	})

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.WaitTimeout(30*time.Second) && token.Error() != nil {
		log.Fatalf("connect failed: %v", token.Error())
	}
	log.Printf("connected: %v", mqttcli.Broker)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGABRT)

	<-quit
}
