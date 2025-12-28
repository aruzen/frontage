package repository

import (
	"encoding/binary"
	"encoding/json"
	"frontage/pkg/network"
	"github.com/google/uuid"
	"log/slog"
	"net"
	"sync"
)

var (
	connectionTable   = make(map[uuid.UUID]*ConnWithMutex)
	connectionTableMu sync.RWMutex
)

type ConnWithMutex struct {
	Conn net.Conn
	Mtx  sync.Mutex
}

func GetConnection(id uuid.UUID) *ConnWithMutex {
	connectionTableMu.RLock()
	defer connectionTableMu.RUnlock()
	connection, ok := connectionTable[id]
	if !ok {
		slog.Error("GetConnection: connection not found")
	}
	return connection
}

func AddConnection(id uuid.UUID, conn net.Conn) {
	connectionTableMu.Lock()
	defer connectionTableMu.Unlock()
	connectionTable[id] = &ConnWithMutex{conn, sync.Mutex{}}
}

func RemoveConnection(id uuid.UUID) {
	connectionTableMu.Lock()
	defer connectionTableMu.Unlock()
	delete(connectionTable, id)
}

func SendPacket(id uuid.UUID, h network.Packet) bool {
	connection := GetConnection(id)
	if connection == nil {
		return false
	}
	body, err := json.Marshal(h)
	if err != nil {
		return false
	}
	header := make([]byte, 6)
	binary.LittleEndian.PutUint16(header[:2], uint16(h.PacketTag()))
	binary.LittleEndian.PutUint32(header[2:6], uint32(len(body)))
	connection.Mtx.Lock()
	defer connection.Mtx.Unlock()
	writeAll := func(data []byte) bool {
		for len(data) > 0 {
			written, err := connection.Conn.Write(data)
			if err != nil {
				return false
			}
			if written <= 0 {
				return false
			}
			data = data[written:]
		}
		return true
	}
	if !writeAll(header) {
		return false
	}
	if !writeAll(body) {
		return false
	}
	return true
}
