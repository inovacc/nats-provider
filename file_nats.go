package natsprovider

import "github.com/nats-io/nats.go"

type FileNATSProvider struct {
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
	file        FileProvider
}

func NewFileNATSProvider(url string) (FileProvider, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	p := &FileNATSProvider{
		name:        "file_nats",
		version:     "1.0.0",
		description: "File NATS Provider",
		nc:          nc,
		js:          js,
	}

	p.file = p

	return p, nil
}

func (p *FileNATSProvider) GetName() string        { return p.name }
func (p *FileNATSProvider) GetVersion() string     { return p.version }
func (p *FileNATSProvider) GetDescription() string { return p.description }

func (p *FileNATSProvider) GetConfig() map[string]any {
	return map[string]any{"url": p.nc.ConnectedUrl()}
}

func (p *FileNATSProvider) Core() CoreProvider {
	return p.core
}

func (p *FileNATSProvider) KeyValue() (KeyValueProvider, error) {
	return p.kv, nil
}

func (p *FileNATSProvider) GetFile(name string) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (p *FileNATSProvider) PutFile(name string, data []byte) error {
	// TODO implement me
	panic("implement me")
}

func (p *FileNATSProvider) DeleteFile(name string) error {
	// TODO implement me
	panic("implement me")
}

func (p *FileNATSProvider) ListFiles() ([]string, error) {
	// TODO implement me
	panic("implement me")
}

func (p *FileNATSProvider) WatchFile(name string, cb func(string, []byte)) error {
	// TODO implement me
	panic("implement me")
}

func (p *FileNATSProvider) UnwatchFile(name string) error {
	// TODO implement me
	panic("implement me")
}

func (p *FileNATSProvider) Close() error {
	// TODO implement me
	panic("implement me")
}
