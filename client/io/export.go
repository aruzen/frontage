package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"sync"

	"frontage/pkg/network"
	"frontage/pkg/network/data"
	"frontage/pkg/network/lobby_api"
)

var (
	clientMu   sync.Mutex
	clientConn net.Conn
)

func writeAll(conn net.Conn, data []byte) error {
	for len(data) > 0 {
		n, err := conn.Write(data)
		if err != nil {
			return err
		}
		if n <= 0 {
			return errors.New("write returned 0 bytes")
		}
		data = data[n:]
	}
	return nil
}

//export InitClient
func InitClient(addr *C.char) C.int {
	if addr == nil {
		return -1
	}
	goAddr := C.GoString(addr)
	conn, err := net.Dial("tcp", goAddr)
	if err != nil {
		return -2
	}
	clientMu.Lock()
	defer clientMu.Unlock()
	if clientConn != nil {
		_ = clientConn.Close()
	}
	clientConn = conn
	return 0
}

//export CloseClient
func CloseClient() C.int {
	clientMu.Lock()
	defer clientMu.Unlock()
	if clientConn != nil {
		_ = clientConn.Close()
		clientConn = nil
	}
	return 0
}

//export SendWaitMatchMake
func SendWaitMatchMake(matchType C.int) C.int {
	clientMu.Lock()
	conn := clientConn
	clientMu.Unlock()
	if conn == nil {
		return -3
	}
	packet := lobby_api.WaitMatchMakePacket{Type: data.MatchType(matchType)}
	body, err := json.Marshal(packet)
	if err != nil {
		return -4
	}
	header := make([]byte, 6)
	binary.LittleEndian.PutUint16(header[:2], uint16(network.WAIT_MATCH_MAKE_PACKET_TAG))
	binary.LittleEndian.PutUint32(header[2:6], uint32(len(body)))
	if err := writeAll(conn, header); err != nil {
		return -5
	}
	if err := writeAll(conn, body); err != nil {
		return -6
	}
	return 0
}

func main() {}
