package saicmqtt

// the topics to subscribe to
type Topic int

const (
	TopicSOC Topic = iota
	TopicIsCharging
	TopicRange
	TopicMileageSinceLastCharge
	TopicMileageOfTheDay
	TopicCurentJourney

	// temperature
	TopicIntTemp
	TopicExtTemp

	// doors/windows
	TopicBoot
	TopicBonnet
	TopicDoors
)

var topicPaths = map[Topic]string{
	TopicSOC:                    "drivetrain/soc",
	TopicIsCharging:             "drivetrain/charging",
	TopicRange:                  "drivetrain/range",
	TopicMileageOfTheDay:        "drivetrain/mileageOfTheDay",
	TopicMileageSinceLastCharge: "drivetrain/mileageSinceLastCharge",
	TopicCurentJourney:          "drivetrain/currentJourney",
	TopicIntTemp:                "climate/interiorTemperature",
	TopicExtTemp:                "climate/exteriorTemperature",
	TopicDoors:                  "doors/locked",
	TopicBoot:                   "doors/boot",
	TopicBonnet:                 "doors/bonnet",
}

var topicNames = map[string]Topic{
	"drivetrain/soc":                    TopicSOC,
	"drivetrain/charging":               TopicIsCharging,
	"drivetrain/range":                  TopicRange,
	"drivetrain/mileageOfTheDay":        TopicMileageOfTheDay,
	"drivetrain/mileageSinceLastCharge": TopicMileageSinceLastCharge,
	"drivetrain/currentJourney":         TopicCurentJourney,
	"climate/interiorTemperature":       TopicIntTemp,
	"climate/exteriorTemperature":       TopicExtTemp,
	"doors/locked":                      TopicDoors,
	"doors/boot":                        TopicBoot,
	"doors/bonnet":                      TopicBonnet,
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
