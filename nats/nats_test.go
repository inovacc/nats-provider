package nats

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestWatchAndSync(t *testing.T) {
	kv, err := testObj.js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket:      "TEST_BUCKET",
		Description: "test bucket for testing purpose",
		Storage:     nats.FileStorage,
		Compression: true,
	})
	if err != nil {
		t.Fatalf("Error creating key-value store: %v", err)
	}

	if err := WatchAndSync(testObj.ctx, kv, "", func(key string, val []byte) {
		t.Logf("Detected update - key: %s, value: %s", key, string(val))
	}); err != nil {
		t.Fatalf("Error setting up WatchAndSync: %v", err)
	}

	key := uuid.NewString()
	if err = SafeWrite(kv, key, func(current []byte) ([]byte, error) {
		return []byte("hello world"), nil
	}); err != nil {
		t.Fatalf("Error writing to key: %v", err)
	}

	// Allow some time for WatchAndSync to pick up the change
	time.Sleep(200 * time.Millisecond)
}
