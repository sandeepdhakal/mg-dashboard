package saicmqtt

import (
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strings"
)

// the prefix used by the MQTT broker
const MQTT_PREFIX string = "saic"

var prefix string

type BrokerInfo struct {
	uri      string
	port     int
	saicUser string
	Url      string
}

func NewBrokerInfo(uri string, port int, saicUser string) *BrokerInfo {
	b := BrokerInfo{uri: uri, port: port, saicUser: saicUser}
	b.Url = fmt.Sprintf("%s:%d", b.uri, b.port)
	return &b
}

type ClientInfo struct {
	user, pass string
	clientId   string
}

func NewClientInfo(user string, pass string, clientId string) *ClientInfo {
	c := ClientInfo{user: user, pass: pass, clientId: clientId}
	return &c
}

type SaicMessageHandler func(Topic, string, string)

type SaicMqttMessage struct {
	Topic Topic
	Msg   string
	Vin   string
}

func ParseMessage(msg mqtt.Message) (SaicMqttMessage, error) {
	parts := strings.SplitAfterN(msg.Topic(), "/", 5)
	topicPath := parts[4]

	if topic, ok := GetTopicFromPath(topicPath); ok {
		return SaicMqttMessage{topic, string(msg.Payload()), ""}, nil
	} else {
		fmt.Println("unrecognised topic")
		return SaicMqttMessage{}, errors.New("unrecognised topic")
	}
}

type SaicMqttClient struct {
	brokerInfo BrokerInfo
	clientInfo ClientInfo
	client     mqtt.Client
}

func NewSaicMqttClient(brokerInfo BrokerInfo,
	clientInfo ClientInfo,
	connectHandler mqtt.OnConnectHandler,
	disconnectHandler mqtt.ConnectionLostHandler,
	messagePubHandler mqtt.MessageHandler,
) SaicMqttClient {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerInfo.Url)
	opts.SetClientID(clientInfo.clientId)
	opts.SetUsername(clientInfo.user)
	opts.SetPassword(clientInfo.pass)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = disconnectHandler
	opts.SetDefaultPublishHandler(messagePubHandler)
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	vin := "+"
	prefix = fmt.Sprintf(
		"%s/%s/vehicles/%s",
		MQTT_PREFIX, brokerInfo.saicUser, vin)

	client := SaicMqttClient{brokerInfo, clientInfo, c}

	return client
}

func (c *SaicMqttClient) Subscribe(topic Topic) {
	topicPath := fmt.Sprintf("%s/%s", prefix, topicPaths[topic])
	token := c.client.Subscribe(topicPath, 1, nil)
	token.Wait()
	// fmt.Printf("Subscribed to %s\n", topicPath)
}

func (c *SaicMqttClient) Disconnect(quiesce uint) {
	c.client.Disconnect(quiesce)
}
