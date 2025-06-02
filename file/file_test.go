package file

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func TestMain(m *testing.M) {
	ns, err := server.NewServer(&server.Options{
		JetStream: true,
		Debug:     true,
		Trace:     true,
	})
	if err != nil {
		log.Fatalf("Error creating nats server: %v", err)
	}

	go ns.Start()

	if !ns.ReadyForConnections(4 * time.Second) {
		log.Fatalf("Error starting nats server: %v", err)
	}

	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		log.Fatal(err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		log.Fatal(err)
	}
	stream := ""
	subject := ""

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     stream,
		Subjects: []string{subject},
		Storage:  nats.FileStorage,
	})
	if err != nil && !errors.Is(err, nats.ErrStreamNameAlreadyInUse) {
		nc.Close()
		log.Fatal(err)
	}
}
