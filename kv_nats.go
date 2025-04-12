package natsprovider

import (
	"errors"
	"sync"

	"github.com/nats-io/nats.go"
)

type kvProvider struct {
	js        nats.JetStreamContext
	store     nats.KeyValue
	storeName string
	watchers  map[string]nats.KeyWatcher
	lock      sync.Mutex
}

func NewKeyValueProvider(js nats.JetStreamContext, storeName string) (KeyValueProvider, error) {
	store, err := js.KeyValue(storeName)
	if errors.Is(err, nats.ErrBucketNotFound) {
		store, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: storeName,
		})
	}
	if err != nil {
		return nil, err
	}
	return &kvProvider{
		js:        js,
		store:     store,
		storeName: storeName,
		watchers:  make(map[string]nats.KeyWatcher),
	}, nil
}

func (kv *kvProvider) Get(key string) (string, error) {
	e, err := kv.store.Get(key)
	if err != nil {
		return "", err
	}
	return string(e.Value()), nil
}

func (kv *kvProvider) Set(key, value string) error {
	_, err := kv.store.PutString(key, value)
	return err
}

func (kv *kvProvider) Delete(key string) error {
	return kv.store.Delete(key)
}

func (kv *kvProvider) List() ([]string, error) {
	keys, err := kv.store.Keys()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (kv *kvProvider) Exists(key string) (bool, error) {
	_, err := kv.store.Get(key)
	if err != nil {
		if errors.Is(err, nats.ErrKeyNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (kv *kvProvider) Close() error {
	return nil // NATS KV does not require explicit close
}

func (kv *kvProvider) Open() error {
	return nil // Already opened during init
}

func (kv *kvProvider) Create() error {
	_, err := kv.js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: kv.storeName,
	})
	return err
}

func (kv *kvProvider) DeleteStore() error {
	return kv.js.DeleteKeyValue(kv.storeName)
}

func (kv *kvProvider) GetStoreName() string {
	return kv.storeName
}

func (kv *kvProvider) Watch(key string, callback func(string, string)) error {
	kv.lock.Lock()
	defer kv.lock.Unlock()

	if _, ok := kv.watchers[key]; ok {
		return nil // already watching
	}

	watcher, err := kv.store.Watch(key)
	if err != nil {
		return err
	}

	kv.watchers[key] = watcher

	go func() {
		for update := range watcher.Updates() {
			if update != nil && update.Operation() != nats.KeyValueDelete {
				callback(update.Key(), string(update.Value()))
			}
		}
	}()

	return nil
}

func (kv *kvProvider) Unwatch(key string) error {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	if w, ok := kv.watchers[key]; ok {
		w.Stop()
		delete(kv.watchers, key)
	}
	return nil
}
