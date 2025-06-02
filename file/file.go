package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/nats-io/nats.go"
)

type JetFile struct {
	js       nats.JetStreamContext
	stream   string
	subject  []string
	readSeq  uint64 // sequence index for reads
	readLast uint64 // upper bound for reading
}

func OpenJetFile(js nats.JetStreamContext, cfg *nats.StreamConfig) (*JetFile, error) {
	return &JetFile{
		js:      js,
		stream:  cfg.Name,
		subject: cfg.Subjects,
	}, nil
}

func (f *JetFile) Write(p []byte) (n int, err error) {
	ack, err := f.js.Publish(f.subject[0], p)
	if err != nil {
		return 0, fmt.Errorf("publish: %w", err)
	}
	log.Printf("published message: seq %d", ack.Sequence)
	return len(p), nil
}

func (f *JetFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.readSeq = uint64(offset)
	case io.SeekCurrent:
		f.readSeq += uint64(offset)
	case io.SeekEnd:
		if offset > 0 {
			return 0, fmt.Errorf("invalid offset beyond end")
		}
		f.readSeq = f.readLast + uint64(offset)
	default:
		return 0, fmt.Errorf("invalid whence")
	}
	return int64(f.readSeq), nil
}

func (f *JetFile) Read(p []byte) (int, error) {
	msg, err := f.js.GetMsg(f.stream, f.readSeq)
	if err != nil {
		return 0, fmt.Errorf("get msg: %w", err)
	}
	n := copy(p, msg.Data)
	f.readSeq++
	return n, nil
}

func (f *JetFile) Close() error {
	// f.conn.Drain()
	// f.conn.Close()
	return nil
}

func (f *JetFile) LoadAll(ctx context.Context) ([]byte, error) {
	sub, err := f.js.PullSubscribe(f.subject[0], "reader")
	if err != nil {
		return nil, fmt.Errorf("subscribe: %w", err)
	}

	var content []byte
	for {
		msgs, err := sub.Fetch(1, nats.Context(ctx))
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				break // done reading
			}
			return nil, fmt.Errorf("fetch: %w", err)
		}
		if len(msgs) == 0 {
			break
		}
		for _, msg := range msgs {
			content = append(content, msg.Data...)
			msg.Ack()
		}
	}
	return content, nil
}
