package actor

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type RabbitmqConfig struct {
	ConnString string
}

//ActorSystem is responsible to manage and relay messages to actors
type ActorSystem struct {
	actors         map[string]*Actor
	connection     *amqp.Connection
	rabbitmqConfig RabbitmqConfig
	started        bool
}

// System creates a new ActorSystem
func System(config RabbitmqConfig) (*ActorSystem, error) {
	conn, err := amqp.Dial(config.ConnString)
	if err != nil {
		logrus.Fatalf("Error while Creating Actor System %v", err)
		return nil, fmt.Errorf("Error Starting Actor System %v", err)
	}
	return &ActorSystem{rabbitmqConfig: config, connection: conn, actors: make(map[string]*Actor), started: false}, nil
}

// Register adds new Actor to the system
func (as *ActorSystem) Register(config ActorConfig, handleFunc HandlerFunc) *Actor {
	publishChannel, err := as.connection.Channel()
	if err != nil {
		logrus.Errorf("Failed to open channel %v", err)
		// return nil, fmt.Errorf("Failed to open mq channel %v", err)
	}
	var subscribeChannels []*ActorHandler
	for i := 0; i < config.Concurrency; i++ {
		ch, err := as.connection.Channel()
		if err != nil {
			logrus.Errorf("Failed to open channel %v", err)
			// return nil, fmt.Errorf("Failed to open mq channel %v", err)
		}
		subscribeChannels = append(subscribeChannels, &ActorHandler{subscribeChannel: ch, handler: handleFunc, name: config.Name})
	}
	actor := Actor{publishChannel: publishChannel, subscribeChannels: subscribeChannels, config: config}
	as.actors[config.Name] = &actor
	return &actor
}

//SendMessage routes messages to actor by name
func (as *ActorSystem) SendMessage(name string, message []byte) error {
	if !as.started {
		return fmt.Errorf("Actor System not initialized")
	}
	if actor, ok := as.actors[name]; ok {
		return actor.SendMessage(message)
	}
	return fmt.Errorf("No Actor with name %s registered", name)
}

//Init starts all actors in the system
func (as *ActorSystem) Init() {
	as.started = true
	for _, v := range as.actors {
		v.Start()
	}
}

//Stop cleans up all actors relenquishing all resources
func (as *ActorSystem) Stop() {
	as.started = false
	for _, v := range as.actors {
		v.Stop()
	}
}
