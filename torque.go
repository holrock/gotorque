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
	"unsafe"
)

// torque connection.
type Torque struct {
	serverName string // conneted server nama.e
	serverID   int    // connection descriptor
}

// GetLastError return last error.
func GetLastError() error {
	en := C.pbs_errno
	estr := C.pbs_strerror(en)
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
func (t *Torque) StatServer() (attr map[string]string, err error) {
	bs := C.pbs_statserver(C.int(t.serverID), nil, nil)
	if bs == nil {
		return nil, GetLastError()
	}
	defer C.pbs_statfree(bs)

	attr = make(map[string]string)
	attr["name"] = C.GoString(bs.name)

	for at := bs.attribs; at != nil; at = at.next {
		attr[C.GoString(at.name)] = C.GoString(at.value)
	}
	return attr, nil
}
