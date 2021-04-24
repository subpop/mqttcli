package `mqttcli` is a program that provides two subcommands (`pub` and `sub`)
that allow command-line level access to an MQTT broker.

`sub` subscribes to a topic and prints messages received to standard output.
`pub` publishes the provided message to the provided topic. Both programs accept
flags that can be provided as a config file.

## Examples

### Flags ###

* `go run ./ -broker tcp://test.mosquitto.org:1883 -topic mqttcli/test sub`
* `go run ./ -broker tcp://test.mosquitto.org:1883 -topic mqttcli/test pub -message "hello"`

### Config File ###

```
cat > sub.cfg << EOF
broker tcp://test.mosquitto.org:1883
topic mqttcli/test
EOF
go run ./ -config sub.cfg
```

```
cat > pub.cfg << EOF
broker tcp://test.mosquitto.org:1883
topic mqttcli/test
EOF
go run ./ -config sub.cfg pub -message test
```
