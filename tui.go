package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"os"
	"os/signal"
	"syscall"

	saic "github.com/sandeepdhakal/mg-dashboard/saicmqtt"
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var disconnectHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v\n", err)
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	if saicMsg, err := saic.ParseMessage(msg); err == nil {
		switch saicMsg.Topic {
		case saic.TopicSOC:
			fmt.Printf("Battery charge: %s %%\n", msg.Payload())
		case saic.TopicIsCharging:
			fmt.Printf("Is charging: %s\n", msg.Payload())
		case saic.TopicRange:
			fmt.Printf("Range: %s kms\n", msg.Payload())
		case saic.TopicIntTemp:
			fmt.Printf("Interior temperature: %s \u00B0C\n", msg.Payload())
		case saic.TopicExtTemp:
			fmt.Printf("Exterior temperature: %s \u00B0C\n", msg.Payload())
		}
	} else {
		fmt.Print("Error parsing message!!")
	}
}

func main() {
	keepAlive := make(chan os.Signal, 1)
	signal.Notify(keepAlive, syscall.SIGINT, syscall.SIGTERM)

	saicUser := "sandeep.dhakal@gmail.com"
	brokerInfo := saic.NewBrokerInfo(
		"tcp://localhost",
		1883,
		saicUser)
	clientInfo := saic.NewClientInfo(
		"mqtt_user",
		"secret",
		"")

	client := saic.NewSaicMqttClient(*brokerInfo,
		*clientInfo,
		connectHandler,
		disconnectHandler,
		messagePubHandler)

	// subscribe to all available topics
	client.Subscribe(saic.TopicSOC)
	client.Subscribe(saic.TopicRange)
	client.Subscribe(saic.TopicIntTemp)
	client.Subscribe(saic.TopicExtTemp)

	<-keepAlive
	client.Disconnect(250)
}
