package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	saic "github.com/sandeepdhakal/mg-dashboard/saicmqtt"
)

// func handleError(e error) {
// 	if e != nil {
// 		fmt.Println("Error:", e)
// 		os.Exit(1)
// 	}
// }

// func parseFloat(s string) float64 {
// 	res, e := strconv.ParseFloat(s, 64)
// 	handleError(e)
// 	return res
// }

type vehicleInfo struct {
	soc     float64
	rng     float64
	intTemp int
	extTemp int
}

var v vehicleInfo = vehicleInfo{}

// model to be used with bubbletea
type model struct {
	sub      chan vehicleInfo
	progress progress.Model
	v        vehicleInfo
}

var m model = model{}

// handler for when connection is established with the mqtt broker
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

// when disconnected from the mqtt broker
var disconnectHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v\n", err)
}

// when a message is received for a subscribed topic
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	if saicMsg, err := saic.ParseMessage(msg); err == nil {
		switch saicMsg.Topic {
		case saic.TopicSOC:
			v.soc, _ = strconv.ParseFloat(string(msg.Payload()), 64)
			m.sub <- v
		case saic.TopicIsCharging:
			fmt.Printf("Is charging: %s\n", msg.Payload())
		case saic.TopicRange:
			v.rng, _ = strconv.ParseFloat(string(msg.Payload()), 64)
			m.sub <- v
		case saic.TopicIntTemp:
			v.intTemp, _ = strconv.Atoi(string(msg.Payload()))
			m.sub <- v
		case saic.TopicExtTemp:
			v.extTemp, _ = strconv.Atoi(string(msg.Payload()))
			m.sub <- v
		}
	} else {
		fmt.Print("Error parsing message!!")
	}
}

// new saic mqtt client
func newClient() saic.SaicMqttClient {
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

	return client
}

// message passed to bubbletea
type updateMsg vehicleInfo

// wait for activity on the model's channel
func waitForActivity(sub chan vehicleInfo) tea.Cmd {
	return func() tea.Msg {
		return updateMsg(<-sub)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForActivity(m.sub), // wait for activity
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case updateMsg:
		m.v = vehicleInfo(msg)
		return m, waitForActivity(m.sub) // wait for next event
	}
	return m, nil
}

func (m model) View() string {
	s := "Battery status\n"
	s += "---------------\n"
	s += m.progress.ViewAs(m.v.soc/100) + "\n"
	s += fmt.Sprintf("Range: %.2f kms\n", m.v.rng)

	s += "\nTemperature\n"
	s += "-----------\n"
	s += fmt.Sprintf("Interior temperature: %d \u00B0C\n", m.v.intTemp)
	s += fmt.Sprintf("Exterior temperature: %d \u00B0C\n", m.v.extTemp)
	s += "\nPress q to quit.\n"
	return s
}

func main() {
	prog := progress.New()
	prog.Width = 30
	m = model{
		v:        vehicleInfo{},
		progress: prog,
		sub:      make(chan vehicleInfo),
	}

	client := newClient()

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}

	client.Disconnect(250)
}
