package nats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

// SafeWrite tries to write data into NATS KV with CAS retry logic.
func SafeWrite(kv nats.KeyValue, key string, modifyFn func(current []byte) ([]byte, error)) error {
	for attempt := 0; attempt < 3; attempt++ {
		entry, err := kv.Get(key)
		switch {
		case errors.Is(err, nats.ErrKeyNotFound):
			newData, modErr := modifyFn(nil)
			if modErr != nil {
				return modErr
			}
			_, err = kv.Create(key, newData)
			if err == nil {
				return nil
			}
			if errors.Is(err, nats.ErrKeyExists) {
				time.Sleep(50 * time.Millisecond)
				continue
			}
			return fmt.Errorf("failed to create key %q: %w", key, err)

		case err != nil:
			return fmt.Errorf("failed to get key %q: %w", key, err)
		}

		newData, modErr := modifyFn(entry.Value())
		if modErr != nil {
			return modErr
		}

		if _, err = kv.Update(key, newData, entry.Revision()); err == nil {
			return nil
		}
		if errors.Is(err, nats.ErrKeyExists) {
			time.Sleep(50 * time.Millisecond)
			continue
		}
		return fmt.Errorf("failed to update key %q: %w", key, err)
	}
	return fmt.Errorf("max retries reached for key %q", key)
}

// WatchAndSync watches a KV bucket prefix and runs syncFn on each update. Cancellable with ctx.
func WatchAndSync(ctx context.Context, kv nats.KeyValue, prefix string, syncFn func(key string, val []byte)) error {
	watcher, err := kv.Watch(fmt.Sprintf("%s>", prefix))
	if err != nil {
		return fmt.Errorf("unable to start KV watch: %w", err)
	}

	go func() {
		defer watcher.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case update, ok := <-watcher.Updates():
				if !ok || update == nil || update.Value() == nil || update.Operation() == nats.KeyValueDelete {
					continue
				}
				syncFn(update.Key(), update.Value())
			}
		}
	}()

	return nil
}

// WatchAndSyncTyped adds JSON decoding of values into a provided type.
func WatchAndSyncTyped[T any](ctx context.Context, kv nats.KeyValue, prefix string, syncFn func(key string, t T)) error {
	return WatchAndSync(ctx, kv, prefix, func(key string, val []byte) {
		var typed T
		if err := json.Unmarshal(val, &typed); err == nil {
			syncFn(key, typed)
		}
	})
}
