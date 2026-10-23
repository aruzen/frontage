package controller

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"frontage/pkg/network"
	"frontage/pkg/network/repository"
	"io"
	"net"
	"time"
)

var (
	LoadHeaderFailedErr = errors.New("load header failed")
	ErrPacketTooLarge   = errors.New("packet body exceeds the receive limit")
)

const HeaderLength = 6

type ReceiveOptions struct {
	MaxPacketBodySize uint32
	// FrameReadTimeout starts at the first header byte and is not extended by
	// partial reads. It also bounds the TLS handshake, but allows idle connections.
	FrameReadTimeout time.Duration
}

func DefaultReceiveOptions() ReceiveOptions {
	return ReceiveOptions{MaxPacketBodySize: 1 << 20, FrameReadTimeout: 30 * time.Second}
}

func ReceiveLoop(ctx context.Context, conn net.Conn, systemChan chan network.UnsolvedPacket,
	lobbyChan chan network.UnsolvedPacket, gameChan *repository.BarrierGameChannel) error {
	return ReceiveLoopWithOptions(ctx, conn, systemChan, lobbyChan, gameChan, DefaultReceiveOptions())
}

func ReceiveLoopWithOptions(ctx context.Context, conn net.Conn, systemChan chan network.UnsolvedPacket,
	lobbyChan chan network.UnsolvedPacket, gameChan *repository.BarrierGameChannel, options ReceiveOptions) error {
	defer conn.Close()
	defaults := DefaultReceiveOptions()
	if options.MaxPacketBodySize == 0 {
		options.MaxPacketBodySize = defaults.MaxPacketBodySize
	}
	if options.FrameReadTimeout < 0 {
		return errors.New("frame read timeout must not be negative")
	}
	if options.FrameReadTimeout == 0 {
		options.FrameReadTimeout = defaults.FrameReadTimeout
	}
	stop := context.AfterFunc(ctx, func() { _ = conn.Close() })
	defer stop()
	if secure, ok := conn.(*tls.Conn); ok {
		handshakeCtx, cancel := context.WithTimeout(ctx, options.FrameReadTimeout)
		err := secure.HandshakeContext(handshakeCtx)
		cancel()
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return err
		}
	}
	readFull := func(data []byte) error {
		_, err := io.ReadFull(conn, data)
		if err != nil && ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
	header := make([]byte, HeaderLength)
	var body []byte
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := conn.SetReadDeadline(time.Time{}); err != nil {
			return err
		}
		if err := readFull(header[:1]); err != nil {
			return err
		}
		if err := conn.SetReadDeadline(time.Now().Add(options.FrameReadTimeout)); err != nil {
			return err
		}
		if err := readFull(header[1:]); err != nil {
			return err
		}
		packetTag := network.PacketTag(binary.LittleEndian.Uint16(header[:2]))
		packetLength := binary.LittleEndian.Uint32(header[2:6])
		if packetLength > options.MaxPacketBodySize || uint64(packetLength) > uint64(^uint(0)>>1) {
			return fmt.Errorf("%w: got %d bytes, limit %d", ErrPacketTooLarge, packetLength, options.MaxPacketBodySize)
		}
		if int(packetLength) > cap(body) {
			body = make([]byte, packetLength)
		} else {
			body = body[:packetLength]
		}
		if err := readFull(body); err != nil {
			return err
		}
		var target chan network.UnsolvedPacket
		switch packetTag >> 14 {
		case network.SystemPacketFlag:
			target = systemChan
		case network.LobbyPacketFlag:
			target = lobbyChan
		case network.GamePacketFlag:
			if gameChan != nil && gameChan.Living.Load() {
				target = gameChan.Chan
			}
		}
		if target == nil {
			continue
		}
		packet := network.UnsolvedPacket{Tag: packetTag, Body: append([]byte(nil), body...)}
		select {
		case target <- packet:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
