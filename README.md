# Go Actors

Package allows creation of Rabbitmq Actors.

# Usage

```go

func main() {
    // Create an Actor system
	system, err := actor.System(actor.RabbitmqConfig{"amqp://localhost"})
	if err != nil {
		fmt.Printf("Error Creating actor System %v", err)
    }
    // Register Actors to system 
	system.Register(actor.ActorConfig{Concurrency: 2, Name: "Example"}, func(message []byte) error {
		var msg Message
		json.Unmarshal(message, &msg)
		fmt.Printf("Message Rx: %v\n", msg)
		return nil
	})
    system.Init()
    
	system.SendMessage("Example", []byte(`{"message": "Hellow WOrld1"}`))

	for {
	}
}

```