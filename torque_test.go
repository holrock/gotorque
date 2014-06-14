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
	server, err := torque.StatServer()
	if err != nil {
		t.Error("cannot get server stat")
	}

	if server.name != torque.ServerName() {
		t.Errorf("server name and hostname unmatched %s %s",
			server.name, torque.ServerName())
	}
}

func TestStatQueue(t *testing.T) {
	torque, err := Connect()
	defer torque.Disconnect()
	queue, err := torque.StatQue()
	if err != nil {
		t.Error("cannot get server stat")
	}
	if len(queue) == 0 {
		t.Error("empty queue infomation")
	}
}

func TestStatNode(t *testing.T) {
	torque, err := Connect()
	defer torque.Disconnect()
	node, err := torque.StatNode()
	if err != nil {
		t.Error("cannot get server stat")
	}
	if len(node) == 0 {
		t.Error("empty node infomation")
	}
}
