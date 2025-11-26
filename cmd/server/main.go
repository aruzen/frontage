package main

import (
	"frontage/pkg/repository"
	"github.com/google/uuid"
	"log/slog"
	"net"
)

func main() {
	for {
		addr, err := net.ResolveUDPAddr("udp", ":8275")
		if err != nil {
			slog.Error("Failed to resolve UDP address", err)
			continue
		}
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			slog.Error("Failed to listen UDP address", err)
			continue
		}
		defer func(conn *net.UDPConn) {
			err := conn.Close()
			if err != nil {
				slog.Error("Failed to close UDP connection", err)
			}
		}(conn)
		id, err := uuid.NewUUID()
		if err != nil {
			slog.Error("Failed to generate UUID", err)
		}
		repository.AddConnection(id, conn)
		defer repository.RemoveConnection(id)
	}

}
