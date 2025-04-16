package broker

import "time"

type NATSConnection struct {
	URL                 string        `mapstructure:"url" validate:"required"`
	ClientName          string        `mapstructure:"client-name" validate:"required"`
	Timeout             time.Duration `mapstructure:"timeout" validate:"required"`
	ReconnectWait       time.Duration `mapstructure:"reconnect-wait" validate:"required"`
	MaxReconnectWait    time.Duration `mapstructure:"max-reconnect-wait" validate:"required"`
	PingInterval        time.Duration `mapstructure:"ping-interval" validate:"required"`
	MaxPingsOutstanding int           `mapstructure:"max-pings-outstanding" validate:"required"`
	Username            string        `mapstructure:"username"`
	Password            string        `mapstructure:"password"`
}

type JetStreamOptions struct {
	MaxWait time.Duration `mapstructure:"max-wait" validate:"required"`
}

type StreamConfig struct {
	Name     string        `mapstructure:"name" validate:"required"`
	Subjects []string      `mapstructure:"subjects" validate:"required"`
	MaxAge   time.Duration `mapstructure:"max-age" validate:"required"`
}

type JetStreamSubscriberConfig struct {
	StreamName string         `mapstructure:"stream-name" validate:"required"`
	Consumer   ConsumerConfig `mapstructure:"consumer" validate:"required"`
}

type ConsumerConfig struct {
	Durable       string        `mapstructure:"durable" validate:"required"`
	Subject       string        `mapstructure:"subject" validate:"required"`
	AckWait       time.Duration `mapstructure:"ack-wait" validate:"required"`
	MaxDeliver    int           `mapstructure:"max-deliver" validate:"required"`
	MaxAckPending int           `mapstructure:"max-ack-pending" validate:"required"`
}

type CoreSubscriberConfig struct {
	Subject string `mapstructure:"subject" validate:"required"`
	Group   string `mapstructure:"queue" validate:"required"`
}
