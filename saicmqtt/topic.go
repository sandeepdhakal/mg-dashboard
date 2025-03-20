package saicmqtt

// the topics to subscribe to
type Topic int

const (
	TopicSOC Topic = iota
	TopicIsCharging
	TopicRange
	TopicIntTemp
	TopicExtTemp
)

var topicPaths = map[Topic]string{
	TopicSOC:        "drivetrain/soc",
	TopicIsCharging: "drivetrain/charging",
	TopicRange:      "drivetrain/range",
	TopicIntTemp:    "climate/interiorTemperature",
	TopicExtTemp:    "climate/exteriorTemperature",
}

var topicNames = map[string]Topic{
	"drivetrain/soc":              TopicSOC,
	"drivetrain/charging":         TopicIsCharging,
	"drivetrain/range":            TopicRange,
	"climate/interiorTemperature": TopicIntTemp,
	"climate/exteriorTemperature": TopicExtTemp,
}

func (t Topic) Path() string {
	return topicPaths[t]
}

func GetTopicFromPath(s string) (Topic, bool) {
	topic, ok := topicNames[s]
	return topic, ok
}
