package natsprovider

import (
	"context"
	"time"

	"github.com/inovacc/nats-provider/utils"
	"github.com/nats-io/nats.go"
)

type coreProvider struct {
	nc *nats.Conn
}

func (c *coreProvider) Publish(subject string, msg []byte, headers map[string]string) error {
	m := &nats.Msg{Subject: subject, Data: msg, Header: nats.Header{}}
	for k, v := range headers {
		m.Header.Set(k, v)
	}
	return c.nc.PublishMsg(m)
}

func (c *coreProvider) Subscribe(subject string, handler MsgHandler) (Unsubscriber, error) {
	sub, err := c.nc.Subscribe(subject, func(m *nats.Msg) {
		handler(&Message{
			Subject: m.Subject,
			Reply:   m.Reply,
			Data:    m.Data,
			Headers: utils.HeaderMap(m.Header),
		})
	})
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (c *coreProvider) QueueSubscribe(subject, queue string, handler MsgHandler) (Unsubscriber, error) {
	sub, err := c.nc.QueueSubscribe(subject, queue, func(m *nats.Msg) {
		handler(&Message{
			Subject: m.Subject,
			Reply:   m.Reply,
			Data:    m.Data,
			Headers: utils.HeaderMap(m.Header),
		})
	})
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (c *coreProvider) Request(subject string, msg []byte, timeoutMs int) (*Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMs)*time.Millisecond)
	defer cancel()
	resp, err := c.nc.RequestWithContext(ctx, subject, msg)
	if err != nil {
		return nil, err
	}
	return &Message{
		Subject: resp.Subject,
		Reply:   resp.Reply,
		Data:    resp.Data,
		Headers: utils.HeaderMap(resp.Header),
	}, nil
}
