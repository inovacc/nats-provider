package natsprovider

import "github.com/nats-io/nats.go"

type (
	Provider interface {
		GetName() string
		GetVersion() string
		GetDescription() string
		GetConfig() map[string]any

		Core() CoreProvider
		KeyValue() (KeyValueProvider, error)
		ObjectStore() (ObjectStoreProvider, error)
		Stream() (StreamProvider, error)
		Config() ConfigProvider // config distribuida con watch
	}

	CoreProvider interface {
		Publish(subject string, msg []byte, headers map[string]string) error
		Subscribe(subject string, handler MsgHandler) (Unsubscriber, error)
		QueueSubscribe(subject, queue string, handler MsgHandler) (Unsubscriber, error)
		Request(subject string, msg []byte, timeoutMs int) (*Message, error)
	}

	MsgHandler func(msg *Message)

	Message struct {
		Subject string
		Reply   string
		Data    []byte
		Headers map[string]string
	}

	KeyValueProvider interface {
		Get(key string) (string, error)
		Set(key, value string) error
		Delete(key string) error
		List() ([]string, error)
		Exists(key string) (bool, error)
		Watch(key string, cb func(string, string)) error
		Unwatch(key string) error
		Close() error
	}

	ConfigProvider interface {
		WatchConfig(key string, cb func(string, string)) error
		GetConfigValue(key string) (string, error)
		SetConfigValue(key, value string) error
	}

	ObjectStoreProvider interface {
		PutObject(name string, data []byte) (*nats.ObjectInfo, error)
		GetObject(name string) ([]byte, error)
		DeleteObject(name string) error
		ListObjects() ([]string, error)
	}

	StreamProvider interface {
		CreateStream(name string, subjects []string) error
		DeleteStream(name string) error
		PublishToStream(stream, subject string, msg []byte, headers map[string]string) error
		SubscribeToStream(stream, durable string, handler MsgHandler) (Unsubscriber, error)

		CreateMirrorStream(name, sourceStream string) error
		CreateSourceStream(name, sourceSubject string) error
	}

	Unsubscriber interface {
		Unsubscribe() error
	}

	FileProvider interface {
		GetFile(name string) ([]byte, error)
		PutFile(name string, data []byte) error
		DeleteFile(name string) error
		ListFiles() ([]string, error)
		WatchFile(name string, cb func(string, []byte)) error
		UnwatchFile(name string) error
		Close() error
	}
)
