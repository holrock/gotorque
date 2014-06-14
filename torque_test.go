package gorque

import (
	"testing"
)

func TestConnect(t *testing.T) {
	torque, err := Connect()
	if err != nil {
		t.Errorf("cannot connect torque server: %s", err)
	}
	torque.Disconnect()
}

func TestStatServer(t *testing.T) {
	torque, err := Connect()
	defer torque.Disconnect()
	server, err := torque.StatServer()
	if err != nil {
		t.Errorf("cannot get server stat: %s", err)
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
		t.Errorf("cannot get queue stat: %s", err)
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
		t.Errorf("cannot get node stat:", err)
	}
	if len(node) == 0 {
		t.Error("empty node infomation")
	}
}

func TestStatJob(t *testing.T) {
	torque, err := Connect()
	defer torque.Disconnect()
	job, err := torque.StatJob("")
	if err != nil {
		t.Errorf("cannot get job stat: ", err)
	}
	if len(job) == 0 {
		t.Error("empty job infomation")
	}
}

func TestEmptyStatJob(t *testing.T) {
	torque, err := Connect()
	defer torque.Disconnect()
	job, err := torque.StatJob("invalidjobid")
	if err != nil && err.Error() != "Unknown queue" {
		t.Errorf("cannot get  stat: %s", err)
	}
	if job != nil {
		t.Error("expect return empty jobs")
	}
}
