package file

import (
	"errors"
	"github.com/nats-io/nats.go"
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	nc, err := nats.Connect(nats.DefaultURL)
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
