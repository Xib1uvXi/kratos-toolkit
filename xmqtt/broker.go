package xmqtt

import (
	"context"
	"errors"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type Handler func(payload []byte) error

type Broker struct {
	client      mqtt.Client
	subscribers *SubscriberSyncMap
	log         *log.Helper
	username    string
	clientID    string
}

func (b *Broker) Start(ctx context.Context) error {
	if b.client.IsConnected() {
		return nil
	}

	if err := b.Connect(); err != nil {
		return err
	}

	return nil
}

func (b *Broker) Stop(ctx context.Context) error {
	if !b.client.IsConnected() {
		return nil
	}

	b.Disconnect()
	return nil
}

func NewBroker(identifier Identifier, logger log.Logger) (*Broker, func(), error) {
	b := &Broker{
		subscribers: NewSubscriberSyncMap(),
		log:         log.NewHelper(log.With(logger)),
	}

	opts := mqtt.NewClientOptions().
		SetCleanSession(true).
		// Because we are dealing with the problem of connection loss by ourself, there is no need for automatic reconnection here
		SetAutoReconnect(false).
		SetResumeSubs(true).
		SetClientID(identifier.ClientID()).
		SetUsername(identifier.Username()).
		SetPassword(identifier.Password()).
		AddBroker(fmt.Sprintf("tcp://%s:1883", identifier.ServerIP())).
		SetConnectionLostHandler(b.onConnectionLost).
		SetOnConnectHandler(b.onConnectHandler)

	b.log.Infof("mqtt target server: %s", fmt.Sprintf("tcp://%s:1883", identifier.ServerIP()))

	client := mqtt.NewClient(opts)
	b.client = client

	b.clientID = identifier.ClientID()
	b.username = identifier.Username()

	if err := b.Connect(); err != nil {
		b.log.Errorf("mqtt connect error: %v", err)
		return nil, nil, err
	}

	clean := func() {
		b.Disconnect()
	}

	return b, clean, nil
}

// Connect
func (b *Broker) Connect() error {
	if b.client.IsConnected() {
		return nil
	}

	token := b.client.Connect()
	if rs, err := checkClientToken(token); !rs {
		return err
	}

	b.log.Infof("mqtt connect success, clientID: %s, username: %s", b.clientID, b.username)

	return nil
}

// Disconnect
func (b *Broker) Disconnect() {
	if b.client.IsConnected() {
		return
	}

	b.client.Disconnect(250)

	b.subscribers.Clear()

	b.log.Infof("mqtt disconnect success")
}

func (b *Broker) Subscribe(topic string, handler func(payload []byte) error) error {
	return b.subscribe(topic, 1, handler)
}

func (b *Broker) Publish(topic string, payload []byte) error {
	return b.publish(topic, payload)
}

func (b *Broker) subscribe(topic string, qos byte, handler Handler) error {
	if !b.client.IsConnected() {
		return errors.New("not connected")
	}

	callback := func(c mqtt.Client, mq mqtt.Message) {
		if err := handler(mq.Payload()); err != nil {
			b.log.Errorf("mqtt message handler error: %v, topic: %s", err, topic)
			return
		}
	}

	if err := b.doSubscribe(topic, qos, callback); err != nil {
		return err
	}

	sub := &Subscriber{
		Topic:    topic,
		Qos:      qos,
		Callback: callback,
	}

	b.subscribers.Add(topic, sub)

	return nil
}

type PubParamsOption func(*PubParams)

type PubParams struct {
	Retained bool
	Qos      byte
}

// Publish
func (b *Broker) publish(topic string, payload []byte, paramOpts ...PubParamsOption) error {
	if !b.client.IsConnected() {
		return errors.New("not connected")
	}

	var param = &PubParams{Retained: false, Qos: 1}

	for _, opt := range paramOpts {
		opt(param)
	}

	token := b.client.Publish(topic, param.Qos, param.Retained, payload)

	return token.Error()
}

// doSubscribe
func (b *Broker) doSubscribe(topic string, qos byte, handler mqtt.MessageHandler) error {
	token := b.client.Subscribe(topic, qos, handler)

	if rs, err := checkClientToken(token); !rs {
		return err
	}

	return nil
}

// onConnectHandler mqtt 连接成功回调
func (b *Broker) onConnectHandler(_ mqtt.Client) {
	b.log.Infof("mqtt on connect")

	b.subscribers.Foreach(func(topic string, sub *Subscriber) {
		if err := b.doSubscribe(sub.Topic, sub.Qos, sub.Callback); err != nil {
			b.log.Error("mqtt broker subscribe message failed:", err)
		}
	})
}

func (b *Broker) onConnectionLost(client mqtt.Client, err error) {
	b.log.Warnf("on connect lost, try to reconnect, error: %v", err)
	b.loopConnect(client)
}

func (b *Broker) loopConnect(client mqtt.Client) {
	for {
		token := client.Connect()
		if rs, err := checkClientToken(token); !rs {
			if err != nil {
				b.log.Errorf("connect error: %s", err.Error())
			}
		} else {
			break
		}
		time.Sleep(5 * time.Second)
	}
}

func checkClientToken(token mqtt.Token) (bool, error) {
	if token.Wait() && token.Error() != nil {
		return false, token.Error()
	}
	return true, nil
}
