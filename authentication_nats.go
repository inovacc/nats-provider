package natsprovider

import "github.com/nats-io/nats.go"

func NewNATSProviderWithAuth(url, username, password, kvStoreName, objStoreName, streamName string) (*NATSProvider, error) {
	nc, err := nats.Connect(url, nats.UserInfo(username, password))
	if err != nil {
		return nil, err
	}
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	core := &coreProvider{nc: nc}
	kv, err := NewKeyValueProvider(js, kvStoreName)
	if err != nil {
		return nil, err
	}
	objStore, err := NewObjectStoreProvider(js, objStoreName)
	if err != nil {
		return nil, err
	}
	stream := NewStreamProvider(js)
	// config := NewConfigProvider(kv)

	p := &NATSProvider{
		name:        "nats",
		version:     "1.0.0",
		description: "NATS and JetStream Provider with Auth",
		nc:          nc,
		js:          js,
		core:        core,
		kv:          kv,
		objStore:    objStore,
		stream:      stream,
		// config:   config,
	}

	return p, nil
}
