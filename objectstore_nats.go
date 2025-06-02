package natsprovider

import (
	"bytes"
	"errors"
	"io"
	"log"

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

func (o *objectStoreProvider) PutObject(name string, data []byte) (*nats.ObjectInfo, error) {
	return o.store.Put(&nats.ObjectMeta{Name: name}, bytes.NewReader(data))
}

func (o *objectStoreProvider) GetObject(name string) ([]byte, error) {
	reader, err := o.store.Get(name)
	if err != nil {
		return nil, err
	}
	defer func(reader nats.ObjectResult) {
		if err := reader.Close(); err != nil {
			log.Printf("Error closing object reader: %v", err)
		}
	}(reader)
	return io.ReadAll(reader)
}

func (o *objectStoreProvider) DeleteObject(name string) error {
	return o.store.Delete(name)
}

func (o *objectStoreProvider) ListObjects() ([]string, error) {
	var names []string
	objects, err := o.store.List()
	if err != nil {
		return nil, err
	}
	for idx := range objects {
		if objects[idx] != nil {
			names = append(names, objects[idx].Name)
		}
	}
	return names, nil
}
