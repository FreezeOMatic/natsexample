package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// docker run -p 4444:4444 nats -p 4444
// default server nats://127.0.0.1:4444

type Message struct {
	Msg   string
	Count int
}

func main() {

	var msg Message

	server := "nats://127.0.0.1:4444"

	nc1, err := nats.Connect(server)
	if err != nil {
		log.Fatal(err)
	}
	nc2, err := nats.Connect(server)
	if err != nil {
		log.Fatal(err)
	}

	defer nc1.Close()
	defer nc2.Close()

	ec1, err := nats.NewEncodedConn(nc1, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer ec1.Close()

	ec2, err := nats.NewEncodedConn(nc2, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer ec2.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		sub(ec1, "testtopic", &msg)

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i <= 10; i++ {
			pub(ec2, "testtopic", Message{Msg: "HI!", Count: i})
			fmt.Println("msg published")
			time.Sleep(time.Second)
		}

	}()
	for {
		time.Sleep(time.Millisecond * 500)
		if msg.Count != 0 {
			fmt.Println(msg)

		}
	}

	//wg.Wait()
}

func sub(ec *nats.EncodedConn, topic string, msg *Message) {

	if _, err := ec.Subscribe(topic, func(msg *Message) {
		log.Printf("Message: %s - Count: %v", msg.Msg, msg.Count)
	}); err != nil {
		log.Fatal(err)

	}
}

func pub(ec *nats.EncodedConn, topic string, message Message) {
	if err := ec.Publish(topic, message); err != nil {
		log.Fatal(err)
	}
}
