package `mqttcli` is a program that provides two subcommands (`pub` and `sub`)
that allow command-line level access to an MQTT broker.

`sub` subscribes to a topic and prints messages received to standard output.
`pub` publishes the provided message to the provided topic. Both programs accept
flags that can be provided as a config file.

## Examples

### Flags ###

* `go run ./cmd/sub -broker tcp://test.mosquitto.org:1883 -topic mqttcli/test`
* `echo hello | go run ./cmd/pub -broker tcp://test.mosquitto.org:1883 -topic mqttcli/test`

### Config File ###

```
cat > sub.cfg << EOF
broker tcp://test.mosquitto.org:1883
topic mqttcli/test
EOF
go run ./cmd/sub -config sub.cfg
```

```
cat > pub.cfg << EOF
broker tcp://test.mosquitto.org:1883
topic mqttcli/test
EOF
echo hello | go run ./cmd/pub -config pub.cfg
```
