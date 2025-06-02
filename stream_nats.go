package natsprovider

import (
	"github.com/inovacc/nats-provider/utils"

	"github.com/nats-io/nats.go"
)

type streamProvider struct {
	js nats.JetStreamContext
}

func NewStreamProvider(js nats.JetStreamContext) StreamProvider {
	return &streamProvider{js: js}
}

func (s *streamProvider) CreateStream(name string, subjects []string) error {
	_, err := s.js.AddStream(&nats.StreamConfig{
		Name:     name,
		Subjects: subjects,
		Storage:  nats.FileStorage,
	})
	return err
}

func (s *streamProvider) DeleteStream(name string) error {
	return s.js.DeleteStream(name)
}

func (s *streamProvider) PublishToStream(stream string, subject string, msg []byte, headers map[string]string) error {
	m := &nats.Msg{Subject: subject, Data: msg, Header: nats.Header{}}
	for k, v := range headers {
		m.Header.Set(k, v)
	}
	_, err := s.js.PublishMsg(m)
	return err
}

func (s *streamProvider) SubscribeToStream(stream string, durableName string, handler MsgHandler) (Unsubscriber, error) {
	sub, err := s.js.PullSubscribe("$JS."+stream+".*", durableName)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			msgs, err := sub.Fetch(1)
			if err != nil {
				continue
			}
			for _, m := range msgs {
				handler(&Message{
					Subject: m.Subject,
					Reply:   m.Reply,
					Data:    m.Data,
					Headers: utils.HeaderMap(m.Header),
				})
				m.Ack()
			}
		}
	}()

	return sub, nil
}

func (s *streamProvider) CreateMirrorStream(name, sourceStream string) error {
	_, err := s.js.AddStream(&nats.StreamConfig{
		Name:   name,
		Mirror: &nats.StreamSource{Name: sourceStream},
	})
	return err
}

func (s *streamProvider) CreateSourceStream(name string, sourceSubject string) error {
	_, err := s.js.AddStream(&nats.StreamConfig{
		Name:     name,
		Subjects: []string{sourceSubject},
	})
	return err
}
