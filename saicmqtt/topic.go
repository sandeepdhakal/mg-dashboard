package saicmqtt

// the topics to subscribe to
type Topic int

const (
	TopicSOC Topic = iota
	TopicIsCharging
	TopicRange
	TopicIntTemp
	TopicExtTemp
	TopicBoot
	TopicBonnet
	TopicDoors
)

var topicPaths = map[Topic]string{
	TopicSOC:        "drivetrain/soc",
	TopicIsCharging: "drivetrain/charging",
	TopicRange:      "drivetrain/range",
	TopicIntTemp:    "climate/interiorTemperature",
	TopicExtTemp:    "climate/exteriorTemperature",
	TopicDoors:      "doors/locked",
	TopicBoot:       "doors/boot",
	TopicBonnet:     "doors/bonnet",
}

var topicNames = map[string]Topic{
	"drivetrain/soc":              TopicSOC,
	"drivetrain/charging":         TopicIsCharging,
	"drivetrain/range":            TopicRange,
	"climate/interiorTemperature": TopicIntTemp,
	"climate/exteriorTemperature": TopicExtTemp,
	"doors/locked":                TopicDoors,
	"doors/boot":                  TopicBoot,
	"doors/bonnet":                TopicBonnet,
}

var BootStatus = [2]string{"Open", "Closed"}
var DoorStatus = [2]string{"Locked", "Unlocked"}

func (t Topic) Path() string {
	return topicPaths[t]
}

func GetTopicFromPath(s string) (Topic, bool) {
	topic, ok := topicNames[s]
	return topic, ok
}
