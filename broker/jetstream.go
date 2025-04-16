package broker

import "github.com/nats-io/nats.go"

type JetStream struct {
	context nats.JetStreamContext
	nc      *nats.Conn
}

func NewJetStream(n *NATS, opts JetStreamOptions) (*JetStream, error) {
	js, err := n.nc.JetStream(nats.MaxWait(opts.MaxWait))
	if err != nil {
		return nil, err
	}

	return &JetStream{
		context: js,
		nc:      n.nc,
	}, nil
}
