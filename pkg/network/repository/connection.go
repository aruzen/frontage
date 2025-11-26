package repository

import (
	"github.com/google/uuid"
	"log/slog"
	"net"
	"sync"
)

var (
	connectionTable   = make(map[uuid.UUID]*net.UDPConn)
	connectionTableMu sync.RWMutex
)

func GetConnection(id uuid.UUID) *net.UDPConn {
	connectionTableMu.RLock()
	defer connectionTableMu.RUnlock()
	connection, ok := connectionTable[id]
	if !ok {
		slog.Error("GetConnection: connection not found")
	}
	return connection
}

func AddConnection(id uuid.UUID, conn *net.UDPConn) {
	connectionTableMu.Lock()
	defer connectionTableMu.Unlock()
	connectionTable[id] = conn
}

func RemoveConnection(id uuid.UUID) {
	connectionTableMu.Lock()
	defer connectionTableMu.Unlock()
	delete(connectionTable, id)
}
