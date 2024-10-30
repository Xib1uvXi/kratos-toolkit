package xmqtt

import "context"

type Identifier interface {
	ClientID() string
	Username() string
	Password() string
	ServerIP() string
}

type MQTT interface {
	Subscribe(topic string, handler Handler) error
	Publish(topic string, payload []byte) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
