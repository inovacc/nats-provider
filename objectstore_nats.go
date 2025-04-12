package natsprovider

import (
	"bytes"
	"errors"
	"io"

	"github.com/nats-io/nats.go"
)

type objectStoreProvider struct {
	store     nats.ObjectStore
	storeName string
}

func NewObjectStoreProvider(js nats.JetStreamContext, storeName string) (ObjectStoreProvider, error) {
	store, err := js.ObjectStore(storeName)
	if errors.Is(err, nats.ErrBucketNotFound) {
		store, err = js.CreateObjectStore(&nats.ObjectStoreConfig{
			Bucket: storeName,
		})
	}
	if err != nil {
		return nil, err
	}
	return &objectStoreProvider{
		store:     store,
		storeName: storeName,
	}, nil
}

func (o *objectStoreProvider) PutObject(name string, data []byte) error {
	return o.store.Put(&nats.ObjectMeta{Name: name}, bytes.NewReader(data))
}

func (o *objectStoreProvider) GetObject(name string) ([]byte, error) {
	reader, err := o.store.Get(name)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func (o *objectStoreProvider) DeleteObject(name string) error {
	return o.store.Delete(name)
}

func (o *objectStoreProvider) ListObjects() ([]string, error) {
	var names []string
	objects := o.store.List()
	for obj := range objects {
		if obj != nil {
			names = append(names, obj.Name)
		}
	}
	return names, nil
}
