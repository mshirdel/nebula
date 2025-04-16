package broker

import (
	"errors"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type NATS struct {
	nc         *nats.Conn
	clientName string
}

func NewNATS(cfg NATSConnection) (*NATS, error) {
	nc, err := nats.Connect(cfg.URL, makeNatsOpts(cfg)...)
	if err != nil {
		return nil, err
	}

	return &NATS{
		nc:         nc,
		clientName: cfg.ClientName,
	}, nil
}

func makeNatsOpts(cfg NATSConnection) []nats.Option {
	opts := []nats.Option{
		nats.Name(cfg.ClientName),
		nats.Timeout(cfg.Timeout),
		nats.ReconnectWait(cfg.ReconnectWait),
		nats.MaxReconnects(int(cfg.MaxReconnectWait / cfg.ReconnectWait)),
		nats.PingInterval(cfg.PingInterval),
		nats.MaxPingsOutstanding(cfg.MaxPingsOutstanding),
		nats.ErrorHandler(natsErrHandler),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			logrus.Infof("nats reconnected to %s", cfg.URL)
		}),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			logrus.Errorf("nats disconnected and will attempt reconnects for %.0fm (%s) error: %v", cfg.MaxReconnectWait.Minutes(), cfg.URL, err)
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			logrus.Errorf("nats closed (%s): %v", cfg.URL, nc.LastError())
		}),
	}

	if cfg.Username != "" && cfg.Password != "" {
		opts = append(opts, nats.UserInfo(cfg.Username, cfg.Password))
	}

	return opts
}

// See https://docs.nats.io/running-a-nats-service/nats_admin/slow_consumers
func natsErrHandler(_ *nats.Conn, sub *nats.Subscription, natsErr error) {
	logrus.Errorf("nats error: %v", natsErr)

	if errors.Is(natsErr, nats.ErrSlowConsumer) {
		pendingMsgs, _, err := sub.Pending()
		if err != nil {
			logrus.Errorf("couldn't get pending messages: %v", err)
			return
		}

		logrus.Infof("falling behind with %d pending messages on subject %q", pendingMsgs, sub.Subject)
	}
}
