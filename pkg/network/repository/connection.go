package repository

import (
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
