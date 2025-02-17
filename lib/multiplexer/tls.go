/*
 * Teleport
 * Copyright (C) 2023  Gravitational, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package multiplexer

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"log/slog"
	"net"
	"time"

	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"
	"golang.org/x/net/http2"

	"github.com/gravitational/teleport"
	"github.com/gravitational/teleport/lib/defaults"
	"github.com/gravitational/teleport/lib/utils"
	logutils "github.com/gravitational/teleport/lib/utils/log"
)

// TLSListenerConfig specifies listener configuration
type TLSListenerConfig struct {
	// Listener is the listener returning *tls.Conn
	// connections on Accept
	Listener net.Listener
	// ID is an identifier used for debugging purposes
	ID string
	// ReadDeadline is a connection read deadline during the TLS handshake (start
	// of the connection). It is set to defaults.HandshakeReadDeadline if
	// unspecified.
	ReadDeadline time.Duration
	// Clock is a clock to override in tests, set to real time clock
	// by default
	Clock clockwork.Clock
}

// CheckAndSetDefaults verifies configuration and sets defaults
func (c *TLSListenerConfig) CheckAndSetDefaults() error {
	if c.Listener == nil {
		return trace.BadParameter("missing parameter Listener")
	}
	if c.ReadDeadline == 0 {
		c.ReadDeadline = defaults.HandshakeReadDeadline
	}
	if c.Clock == nil {
		c.Clock = clockwork.NewRealClock()
	}
	return nil
}

// NewTLSListener returns a new TLS listener
func NewTLSListener(cfg TLSListenerConfig) (*TLSListener, error) {
	if err := cfg.CheckAndSetDefaults(); err != nil {
		return nil, trace.Wrap(err)
	}
	context, cancel := context.WithCancel(context.TODO())
	return &TLSListener{
		log:           slog.With(teleport.ComponentKey, teleport.Component("mxtls", cfg.ID)),
		cfg:           cfg,
		http2Listener: newListener(context, cfg.Listener.Addr()),
		httpListener:  newListener(context, cfg.Listener.Addr()),
		cancel:        cancel,
		context:       context,
	}, nil
}

// TLSListener wraps tls.Listener and detects negotiated protocol
// (assuming it's either http/1.1 or http/2)
// and forwards the appropriate responses to either HTTP/1.1 or HTTP/2
// listeners
type TLSListener struct {
	log           *slog.Logger
	cfg           TLSListenerConfig
	http2Listener *Listener
	httpListener  *Listener
	cancel        context.CancelFunc
	context       context.Context
}

// HTTP2 returns HTTP2 listener
func (l *TLSListener) HTTP2() net.Listener {
	return l.http2Listener
}

// HTTP returns HTTP listener
func (l *TLSListener) HTTP() net.Listener {
	return l.httpListener
}

// Serve accepts and forwards tls.Conn connections
func (l *TLSListener) Serve() error {
	for {
		conn, err := l.cfg.Listener.Accept()
		if err == nil {
			tlsConn, ok := conn.(*tls.Conn)
			if !ok {
				conn.Close()
				l.log.LogAttrs(l.context, slog.LevelError, "Received a non-TLS connection",
					slog.Any("src_addr", logutils.StringerAttr(conn.RemoteAddr())),
					slog.Any("dst_addr", logutils.StringerAttr(conn.LocalAddr())),
					slog.Any("conn_type", logutils.TypeAttr(conn)),
				)
				continue
			}
			go l.detectAndForward(tlsConn)
			continue
		}
		if utils.IsUseOfClosedNetworkError(err) {
			<-l.context.Done()
			return nil
		}
		select {
		case <-l.context.Done():
			return nil
		case <-time.After(5 * time.Second):
		}
	}
}

func (l *TLSListener) detectAndForward(conn *tls.Conn) {
	err := conn.SetReadDeadline(l.cfg.Clock.Now().Add(l.cfg.ReadDeadline))
	if err != nil {
		l.log.LogAttrs(l.context, slog.LevelDebug, "Failed to set connection deadline",
			slog.Any("error", err),
		)
		conn.Close()
		return
	}

	start := l.cfg.Clock.Now()
	if err := conn.HandshakeContext(l.context); err != nil {
		if !errors.Is(trace.Unwrap(err), io.EOF) {
			l.log.LogAttrs(l.context, slog.LevelWarn, "Handshake failed",
				slog.Any("src_addr", logutils.StringerAttr(conn.RemoteAddr())),
				slog.Any("dst_addr", logutils.StringerAttr(conn.LocalAddr())),
				slog.Any("error", err),
			)
		}
		conn.Close()
		return
	}

	// Log warning if TLS handshake takes more than one second to help debug
	// latency issues.
	if elapsed := time.Since(start); elapsed > 1*time.Second {
		l.log.LogAttrs(l.context, slog.LevelWarn, "Slow TLS handshake",
			slog.Any("src_addr", logutils.StringerAttr(conn.RemoteAddr())),
			slog.Any("dst_addr", logutils.StringerAttr(conn.LocalAddr())),
			slog.Duration("handshake_duration", time.Since(start)),
		)
	}

	err = conn.SetReadDeadline(time.Time{})
	if err != nil {
		l.log.WarnContext(l.context, "Failed to reset read deadline", "error", err)
		conn.Close()
		return
	}

	switch conn.ConnectionState().NegotiatedProtocol {
	case http2.NextProtoTLS:
		l.http2Listener.HandleConnection(l.context, conn)
	case teleport.HTTPNextProtoTLS, "":
		l.httpListener.HandleConnection(l.context, conn)
	default:
		conn.Close()
		l.log.LogAttrs(l.context, slog.LevelError, "rejecting connection with unsupported protocol",
			slog.Any("error", err),
			slog.String("protocol", conn.ConnectionState().NegotiatedProtocol),
			slog.Any("src_addr", logutils.StringerAttr(conn.RemoteAddr())),
			slog.Any("dst_addr", logutils.StringerAttr(conn.LocalAddr())),
		)
	}
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *TLSListener) Close() error {
	defer l.cancel()
	return l.cfg.Listener.Close()
}

// Addr returns the listener's network address.
func (l *TLSListener) Addr() net.Addr {
	return l.cfg.Listener.Addr()
}
