package xmqtt

import (
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"sync"
)

type Subscriber struct {
	Topic    string
	Qos      byte
	Callback mqtt.MessageHandler
}

type SubscriberSyncMap struct {
	sync.RWMutex
	subscribers map[string]*Subscriber
}

func NewSubscriberSyncMap() *SubscriberSyncMap {
	return &SubscriberSyncMap{subscribers: make(map[string]*Subscriber)}
}

func (sm *SubscriberSyncMap) Add(topic string, sub *Subscriber) {
	sm.Lock()
	defer sm.Unlock()

	sm.subscribers[topic] = sub
}

func (sm *SubscriberSyncMap) Remove(topic string) error {
	sm.Lock()
	defer sm.Unlock()

	if _, ok := sm.subscribers[topic]; ok {
		delete(sm.subscribers, topic)
		return nil
	} else {
		return errors.New(fmt.Sprintf("topic[%s] not found", topic))
	}
}

func (sm *SubscriberSyncMap) Clear() {
	sm.Lock()
	defer sm.Unlock()

	sm.subscribers = make(map[string]*Subscriber)
}

func (sm *SubscriberSyncMap) Get(topic string) *Subscriber {
	sm.RLock()
	defer sm.RUnlock()

	return sm.subscribers[topic]
}

func (sm *SubscriberSyncMap) Foreach(fnc func(topic string, sub *Subscriber)) {
	sm.RLock()
	defer sm.RUnlock()

	for k, v := range sm.subscribers {
		fnc(k, v)
	}
}
