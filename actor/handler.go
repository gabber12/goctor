package actor

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

//ActorHandler holds a rabbitmq channel to receive messages
type ActorHandler struct {
	name             string
	subscribeChannel *amqp.Channel
	handler          HandlerFunc
}

//Start starts basic consume on the rabbitmq channel
func (h *ActorHandler) Start() {
	shovel, err := h.subscribeChannel.Consume(h.name, h.name, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("basic.consume Failed: %v", err)
	}

	for msg := range shovel {
		err := h.handlerFunc(&msg)
		if err := h.handleAck(msg, err); err != nil {
			logrus.Errorf("Error Acking message %s Actor %s Error %v", msg.Body, h.name, err)
		}
	}
}

func (h *ActorHandler) handlerFunc(message *amqp.Delivery) error {
	logrus.Infof("Actor Name: (%s) Receiving: [%s]\n", h.name, message.Body)
	return h.handler(message.Body)
}

func (h *ActorHandler) handleAck(msg amqp.Delivery, handlerError error) error {
	if handlerError != nil {
		logrus.Errorf("Error running Handler Actor %s Error %v", h.name, handlerError)
		return msg.Nack(false, true)
	}
	return msg.Ack(false)
}

func (h *ActorHandler) stop() error {
	return h.subscribeChannel.Close()
}
