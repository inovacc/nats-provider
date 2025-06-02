package natsprovider

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

var testObj testStruct

type testStruct struct {
	js  nats.JetStreamContext
	ctx context.Context
}

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
		log.Fatalf("Error connecting to nats server: %v", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Error getting JetStream context: %v", err)
	}

	testObj.ctx = ctx
	testObj.js = js

	os.Exit(m.Run())
}
