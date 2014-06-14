// gorque package is a gateway for Torque(OpenPBS) server.
package gorque

/*
#cgo LDFLAGS: -ltorque
#include <stdlib.h>
#include <pbs_ifl.h>
#include <pbs_error.h>
*/
import "C"

import (
	"errors"
	"strings"
	"unsafe"
)

// torque connection.
type Torque struct {
	serverName string // conneted server nama.e
	serverID   int    // connection descriptor
}

// Attribute is wrapper of pbs_ifl attibute.
type Attribute struct {
	name     string
	resource string
	value    string
	op       string
}

// Server data.
type Server struct {
	name string
	attr map[string]Attribute
}

// Queue data.
type Queue struct {
	name string
	attr map[string]Attribute
}

// Node data.
type Node struct {
	name   string
	attr   map[string]Attribute
	status map[string]string
}

// Job data.
type Job struct {
	name         string
	attr         map[string]Attribute
	variableList map[string]string
}

// batchOp is mapping enum to string.
var batchOp = map[C.enum_batch_op]string{
	C.SET:   "SET",
	C.UNSET: "UNSET",
	C.INCR:  "INCR",
	C.DECR:  "DECR",
	C.EQ:    "EQ",
	C.NE:    "NE",
	C.GE:    "GE",
	C.GT:    "GT",
	C.LE:    "LE",
	C.LT:    "ET",
	C.DFLT:  "DFLT",
	C.MERGE: "MERGE",
}

// GetLastError return last error.
func GetLastError() error {
	en := C.pbs_errno
	estr := C.pbs_strerror(en) // static string
	return errors.New(C.GoString(estr))
}

// DefaultServername return pbs_server.
func DefaultServerName() string {
	return C.GoString(C.pbs_server)
}

// ConnectServer open connection to serverName.
func ConnectServer(serverName string) (*Torque, error) {
	srv := C.CString(serverName)
	con := int(C.pbs_connect(srv))
	C.free(unsafe.Pointer(srv))

	if con < 0 {
		return nil, GetLastError()
	}

	return &Torque{
			serverName: serverName,
			serverID:   con},
		nil
}

// Connect open connection to default server.
func Connect() (*Torque, error) {
	srv := DefaultServerName()
	return ConnectServer(srv)
}

// Disconnect close connection.
func (t *Torque) Disconnect() {
	C.pbs_disconnect(C.int(t.serverID))
	t.serverID = -1
}

// ServerName return connected server name.
func (t *Torque) ServerName() string {
	return t.serverName
}

// StatServer return statserver.
func (t *Torque) StatServer() (Server, error) {
	bs := C.pbs_statserver(C.int(t.serverID), nil, nil)
	if bs == nil {
		return Server{}, GetLastError()
	}
	defer C.pbs_statfree(bs)

	srv := Server{name: C.GoString(bs.name)}

	srv.attr = make(map[string]Attribute)
	for name, attr := range attrlToAttributeMap(bs.attribs) {
		srv.attr[name] = attr
	}

	return srv, nil
}

// StatQeueu return stat all queue.
func (t *Torque) StatQue() ([]Queue, error) {
	bs := C.pbs_statque(C.int(t.serverID), nil, nil, nil)
	if bs == nil {
		return nil, GetLastError()
	}
	defer C.pbs_statfree(bs)

	queues := make([]Queue, 0, 1)

	for cur := bs; cur != nil; cur = cur.next {
		q := Queue{}
		q.name = C.GoString(cur.name)
		q.attr = make(map[string]Attribute)

		for name, attr := range attrlToAttributeMap(cur.attribs) {
			q.attr[name] = attr
		}
		queues = append(queues, q)
	}

	return queues, nil
}

// StatNode return stat all nodes.
func (t *Torque) StatNode() ([]Node, error) {
	bs := C.pbs_statnode(C.int(t.serverID), nil, nil, nil)
	if bs == nil {
		return nil, GetLastError()
	}
	defer C.pbs_statfree(bs)

	nodes := make([]Node, 0, 1)

	for cur := bs; cur != nil; cur = cur.next {
		n := Node{}
		n.name = C.GoString(cur.name)
		n.attr = make(map[string]Attribute)

		for name, attr := range attrlToAttributeMap(cur.attribs) {
			n.attr[name] = attr
		}

		status := n.attr["status"]
		delete(n.attr, "status")
		n.status = kvlistToMap(status.value)

		nodes = append(nodes, n)
	}
	return nodes, nil
}

// StatJob return stat all jobs.
func (t *Torque) StatJob() ([]Job, error) {
	bs := C.pbs_statjob(C.int(t.serverID), nil, nil, nil)
	if bs == nil {
		return nil, GetLastError()
	}
	defer C.pbs_statfree(bs)

	jobs := make([]Job, 0, 1)

	for cur := bs; cur != nil; cur = cur.next {
		j := Job{}
		j.name = C.GoString(cur.name)
		j.attr = make(map[string]Attribute)

		for name, attr := range attrlToAttributeMap(cur.attribs) {
			j.attr[name] = attr
		}

		vl := j.attr["Variable_List"]
		delete(j.attr, "Variable_List")
		j.variableList = kvlistToMap(vl.value)

		jobs = append(jobs, j)
	}

	return jobs, nil
}

func attrlToAttributeMap(attrl *C.struct_attrl) map[string]Attribute {
	attrmap := make(map[string]Attribute)

	for attr := attrl; attr != nil; attr = attr.next {
		op := C.enum_batch_op(attr.op)
		sop := batchOp[op]
		name := C.GoString(attr.name)
		key := name
		if attr.resource != nil {
			key += "." + C.GoString(attr.resource)
		}
		attrmap[key] = Attribute{
			name:     name,
			value:    C.GoString(attr.value),
			resource: C.GoString(attr.resource),
			op:       sop}
	}
	return attrmap
}

func kvlistToMap(kvlist string) map[string]string {
	strmap := make(map[string]string)

	for _, s := range strings.Split(kvlist, ",") {
		kv := strings.Split(s, "=")
		if len(kv) == 2 {
			strmap[kv[0]] = kv[1]
		}
	}
	return strmap
}
