package main

import (
	"flag"
	"log"
	"net/http"

	webhooktomqtt "github.com/parkr/webhook-to-mqtt"

	"gosrc.io/mqtt"
)

func main() {
	var mqttServer string
	flag.StringVar(&mqttServer, "mqtt", "tcp://localhost:1883", "MQTT server")
	var httpListen string
	flag.StringVar(&httpListen, "http", ":8080", "HTTP listen address")
	var mqttClientID string
	flag.StringVar(&mqttClientID, "client-id", "webhook-to-mqtt", "MQTT client ID")
	flag.Parse()

	client := mqtt.NewClient(mqttServer)
	client.ClientID = mqttClientID
	client.Messages = make(chan mqtt.Message)
	clientManager := mqtt.NewClientManager(client, nil)
	clientManager.Start()
	defer clientManager.Stop()

	mux := http.NewServeMux()

	httpPrefix := "/publish/"
	mux.Handle(httpPrefix, webhooktomqtt.NewHandler(client, httpPrefix))

	if err := http.ListenAndServe(httpListen, mux); err != nil {
		log.Fatalf("error serving HTTP: %s", err)
	}
}
