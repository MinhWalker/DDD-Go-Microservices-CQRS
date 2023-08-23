package kafka

import (
	"context"
)

// NewKafkaConn create new kafka connection
func NewKafkaConn(ctx context.Context, kafkaCfg *Config) (*kafka.Conn, error) {
	return kafka.DialContext(ctx, "tcp", kafkaCfg.Brokers[0])
}
