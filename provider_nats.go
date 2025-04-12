package natsprovider

import "github.com/nats-io/nats.go"

type NATSProvider struct {
	name        string
	version     string
	description string
	nc          *nats.Conn
	js          nats.JetStreamContext
	core        CoreProvider
	kv          KeyValueProvider
	objStore    ObjectStoreProvider
	stream      StreamProvider
	config      ConfigProvider
}

func NewNATSProvider(url string) (*NATSProvider, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	p := &NATSProvider{
		name:        "nats",
		version:     "1.0.0",
		description: "NATS and JetStream Provider",
		nc:          nc,
		js:          js,
	}

	// Acá inicializarías los subproviders con la conexión nc y js
	// p.core = NewCoreProvider(nc)
	// p.kv = NewKeyValueProvider(js)
	// p.objStore = NewObjectStoreProvider(js)
	// p.stream = NewStreamProvider(js)
	// p.config = NewConfigProvider(p.kv)

	return p, nil
}

func (p *NATSProvider) GetName() string        { return p.name }
func (p *NATSProvider) GetVersion() string     { return p.version }
func (p *NATSProvider) GetDescription() string { return p.description }
func (p *NATSProvider) GetConfig() map[string]any {
	return map[string]any{"url": p.nc.ConnectedUrl()}
}

func (p *NATSProvider) Core() CoreProvider {
	return p.core
}

func (p *NATSProvider) KeyValue() (KeyValueProvider, error) {
	return p.kv, nil
}

func (p *NATSProvider) ObjectStore() (ObjectStoreProvider, error) {
	return p.objStore, nil
}

func (p *NATSProvider) Stream() (StreamProvider, error) {
	return p.stream, nil
}

func (p *NATSProvider) Config() ConfigProvider {
	return p.config
}
