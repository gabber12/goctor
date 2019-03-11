package actor

import (
	"fmt"
	"testing"
)

func TestPublishUnInitializedSystem(t *testing.T) {

	system, err := System(RabbitmqConfig{ConnString: "amqp://localhost"})
	if err != nil {

	}
	system.Register(ActorConfig{Concurrency: 2, Name: "Actor"}, func(message []byte) error {
		fmt.Printf("Rx: %s\n", message)
		return nil
	})

	if err := system.SendMessage("ROLO", []byte("LORO")); err == nil {
		t.Errorf("No Error Raised on Publish on Uninitialized System")

	}

}

func TestPublishInitializedSystem(t *testing.T) {

	system, err := System(RabbitmqConfig{ConnString: "amqp://localhost"})
	if err != nil {

	}
	system.Register(ActorConfig{Concurrency: 2, Name: "Actor"}, func(message []byte) error {
		fmt.Printf("Rx: %s\n", message)
		return nil
	})
	system.Init()

	if err := system.SendMessage("Actor", []byte("LORO")); err != nil {
		t.Errorf("Error Raised on Publish on Uninitialized System")
	}

}
func TestPublishOnUnkownName(t *testing.T) {

	system, err := System(RabbitmqConfig{ConnString: "amqp://localhost"})
	if err != nil {

	}
	system.Register(ActorConfig{Concurrency: 2, Name: "Actor"}, func(message []byte) error {
		fmt.Printf("Rx: %s\n", message)
		return nil
	})
	system.Init()
	if err := system.SendMessage("ROLO", []byte("LORO")); err != nil {
		return
	}
	t.Errorf("Error Raised on Publish on Uninitialized System %v", err)

}

func TestActorRegistration(t *testing.T) {

	system, err := System(RabbitmqConfig{ConnString: "amqp://localhost"})
	if err != nil {

	}
	system.Register(ActorConfig{Concurrency: 2, Name: "Actor"}, func(message []byte) error {
		fmt.Printf("Rx: %s\n", message)
		return nil
	})
	if len(system.actors) != 1 {
		t.Error("Check actor Registration")
	}
	if _, ok := system.actors["Actor"]; !ok {
		t.Error("Actor not registered")
	}
}
