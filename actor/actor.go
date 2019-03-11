package actor

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Actor struct {
	publishChannel    *amqp.Channel
	subscribeChannels []*ActorHandler
	config            ActorConfig
}

type ActorConfig struct {
	Concurrency int
	Name        string
	NumRetries  int
}

type HandlerFunc func(message []byte) error

func (a *Actor) SendMessage(message []byte) error {
	return a.publishChannel.Publish(a.config.Name+"_exchange", a.config.Name, false, false, amqp.Publishing{
		ContentType:  "text/plain",
		Timestamp:    time.Now(),
		DeliveryMode: amqp.Persistent,
		Body:         message,
	})
}

func (a *Actor) Start() {
	a.start()

}

func (a *Actor) Stop() {
	a.publishChannel.Close()
	for _, h := range a.subscribeChannels {
		h.subscribeChannel.Close()
	}
}

func (a *Actor) start() {
	name := a.config.Name
	if err := a.publishChannel.ExchangeDeclare(name+"_exchange", "direct", true, false, false, false, nil); err != nil {
		log.Fatalf("exchange.declare destination: %s", err)
	}

	if _, err := a.publishChannel.QueueDeclare(name, true, false, false, false, nil); err != nil {
		log.Fatalf("queue.declare source: %s", err)
	}

	if err := a.publishChannel.QueueBind(name, name, name+"_exchange", false, nil); err != nil {
		log.Fatalf("queue.bind source: %s", err)
	}

	for _, handler := range a.subscribeChannels {
		go handler.Start()
	}
}
