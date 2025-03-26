package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	saic "github.com/sandeepdhakal/mg-dashboard/saicmqtt"
)

// Vehicle info to be passed to Bubbletea model
type vehicleInfo struct {
	soc                    float64
	rng                    float64
	intTemp                int
	extTemp                int
	bootLocked             bool
	doorsLocked            bool
	mileageSinceLastCharge float64
}

var v vehicleInfo = vehicleInfo{}

// Model to be used with Bubbletea
type model struct {
	sub      chan vehicleInfo
	progress progress.Model
	v        vehicleInfo
}

var m model = model{}

// handler for when connection is established with the mqtt broker
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	// fmt.Println("Connected")
	// TODO: update UI accordingly
}

// when disconnected from the mqtt broker
var disconnectHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	// fmt.Printf("Connection lost: %v\n", err)
	// TODO: update UI accordingly
}

// when a message is received for a subscribed topic
// Here we will pass the updated vehicle object to the model's channel
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
		case saic.TopicBoot:
			v.bootLocked, _ = strconv.ParseBool(string(msg.Payload()))
			m.sub <- v
		case saic.TopicDoors:
			v.doorsLocked, _ = strconv.ParseBool(string(msg.Payload()))
			m.sub <- v
		case saic.TopicMileageSinceLastCharge:
			v.mileageSinceLastCharge, _ = strconv.ParseFloat(string(msg.Payload()), 64)
			m.sub <- v
		}
	} else {
		fmt.Print("Error parsing message!!")
		// TODO: update UI accordingly
	}
}

// Set of methods and struct to read MQTT configuration from environment variables
func getEnvironmentVariable(v string) string {
	value, ok := os.LookupEnv(v)
	if !ok {
		fmt.Printf("Environment variable %s not set.\n", v)
	}
	return value
}

type config struct {
	saicUser   string
	brokerUri  string
	brokerPort int
	brokerUser string
	brokerPass string
}

func getConfig() config {

	// since port is an integer, we need to process it differently
	port, e := strconv.Atoi(getEnvironmentVariable("SAIC_BROKER_PORT"))
	if e != nil {
		fmt.Printf("Environment variable SAIC_BROKER_PORT must be an integer.\n")
		os.Exit(1)
	}

	return config{
		getEnvironmentVariable("SAIC_USER"),
		getEnvironmentVariable("SAIC_BROKER_URI"),
		port,
		getEnvironmentVariable("SAIC_MQTT_USER"),
		getEnvironmentVariable("SAIC_MQTT_PASS"),
	}
}

// new saic mqtt client
func newClient() saic.SaicMqttClient {
	config := getConfig()
	brokerInfo := saic.NewBrokerInfo(config.brokerUri, config.brokerPort, config.saicUser)
	clientInfo := saic.NewClientInfo(config.brokerUser, config.brokerPass, "")

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
	client.Subscribe(saic.TopicDoors)
	client.Subscribe(saic.TopicBoot)
	client.Subscribe(saic.TopicMileageSinceLastCharge)

	return client
}

// message passed to bubbletea
type updateMsg vehicleInfo

// Bubbletea: wait for activity on the model's channel
func waitForActivity(sub chan vehicleInfo) tea.Cmd {
	return func() tea.Msg {
		return updateMsg(<-sub)
	}
}

// Bubbletea: initialisation of the UI
func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForActivity(m.sub), // wait for activity
	)
}

// Bubbletea: update the UI when we receive new vehicle updates or other messages
// such as keypress to interact with the UI
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

func getStringByFlag(flag bool, values [2]string) string {
	if flag {
		return values[0]
	} else {
		return values[1]
	}
}

// Styles
var headingStyle = lipgloss.NewStyle().
	Bold(true).
	Underline(true).
	Height(2)
	// Padding(1, 0, 1, 0)

// Bubbletea: the UI
func (m model) View() string {

	s := headingStyle.Render("Battery Status") + "\n"
	// lipgloss.NewStyle().Border(lipgloss.ThickBorder(), true, false)
	s += m.progress.ViewAs(m.v.soc/100) + "\n"
	s += fmt.Sprintf("Range: %.2f kms\n", m.v.rng)
	s += fmt.Sprintf("Since Last Charge: %.2f kms\n", m.v.mileageSinceLastCharge)

	s += "\n" + headingStyle.Render("Temperature") + "\n"
	s += fmt.Sprintf("Interior temperature: %d \u00B0C\n", m.v.intTemp)
	s += fmt.Sprintf("Exterior temperature: %d \u00B0C\n", m.v.extTemp)

	s += "\n" + headingStyle.Render("Doors/Windows") + "\n"
	s += fmt.Sprintf("Boot: %s\n", getStringByFlag(m.v.bootLocked, saic.BootStatus))
	s += fmt.Sprintf("Doors: %s\n", getStringByFlag(m.v.doorsLocked, saic.DoorStatus))
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
