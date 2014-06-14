package gorque

import (
	"testing"
)

func TestConnect(t *testing.T) {
	torque, err := Connect()
	if err != nil {
		t.Error("cannot connect torque server")
	}
	torque.Disconnect()
}

func TestStatServer(t *testing.T) {
	torque, err := Connect()
	defer torque.Disconnect()
	attr, err := torque.StatServer()
	if err != nil {
		t.Error("cannot get server stat")
	}

	if attr["name"] != torque.ServerName() {
		t.Errorf("server name and hostname unmatched %s %s",
			attr["name"], torque.ServerName())
	}
}
