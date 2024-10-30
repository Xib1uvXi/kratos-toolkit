package xmqtt

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type testIdentify struct {
	client   string
	username string
	password string
}

func (t *testIdentify) ClientID() string {
	return t.client + "-utest"
}

func (t *testIdentify) Username() string {
	return t.username
}

func (t *testIdentify) Password() string {
	return t.password
}

func (t *testIdentify) ServerIP() string {
	//return "159.75.172.107"
	return "127.0.0.1"
}

func TestNewBroker(t *testing.T) {
	t.Skip("local test")
	identifier := &testIdentify{client: "test2", username: "skyline_dev", password: "Hu6GXtkHAr55Qhd"}
	b1, clean1, err := NewBroker(identifier, log.DefaultLogger)
	require.NoError(t, err)
	defer clean1()
	require.NotNil(t, b1)

	data := fmt.Sprintf(`{"action":"format","params":"{\"disk_position\":1}"}`)

	err = b1.Publish("hermes/device/action/aaa", []byte(data))
	require.NoError(t, err)

	time.Sleep(20 * time.Second)
}
