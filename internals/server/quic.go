package server

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"

	"github.com/barnabasSol/retro-rumble/internals/event"
	"github.com/barnabasSol/retro-rumble/internals/hub"
	"github.com/barnabasSol/retro-rumble/internals/models"
	"github.com/google/uuid"
	"github.com/quic-go/quic-go"
)

type QuicServer struct {
	addr string
	lis  *quic.Listener
	hub  *hub.GameHub
}

func NewQuicServer(
	addr string,
	tlsConfig *tls.Config,
	quicConfig *quic.Config,
	lis *quic.Listener,
	hub *hub.GameHub,
) *QuicServer {
	lis, err := quic.ListenAddr(
		addr,
		generateTLSConfig(),
		quicConfig,
	)
	if err != nil {
		panic(err)
	}
	return &QuicServer{
		addr: addr,
		hub:  hub,
		lis:  lis,
	}
}

func (s *QuicServer) Start() error {
	defer s.lis.Close()
	for {
		conn, err := s.lis.Accept(context.Background())
		_ = conn
		if err != nil {
			return err
		}
		p := &models.Player{
			Id:          uuid.NewString(),
			Send:        make(chan []byte, 256),
			QuicStreams: make(models.QuicStream),
		}
		s.hub.Register <- p
		go s.handleConnection(conn, p)
	}
}

func (s *QuicServer) handleConnection(
	conn *quic.Conn,
	p *models.Player,
) {
	defer func() {
		conn.CloseWithError(0, "")
		s.hub.Unregister <- p
		fmt.Println("Server: connection closed of player", p.Id)
	}()
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			fmt.Println("Server: connection closed:", err)
			return
		}

		fmt.Printf(
			"Server: accepted stream %d\n",
			stream.StreamID(),
		)

		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Server: stream read error:", err)
			}
			return
		}

		stream_identifier := string(buf[:n])
		if _, ok := ValidStreams[stream_identifier]; !ok {
			fmt.Printf(
				"Server: invalid stream type '%s'\n",
				stream_identifier,
			)
			return
		}

		json.Marshal(event.NewStreamPayload{
			Identifier: stream_identifier,
			Stream: stream,
		})

		s.hub.GameEvent <- models.InboundEvent{
			Ev: &models.Event{
				Type: event.TypeNewStream,
				Event: ,
			},
			Player: p,
		}

		go func(
			s *quic.Stream,
			p *models.Player,
			hub *hub.GameHub,
		) {
			defer s.Close()
			buf := make([]byte, 1024)
			for {
				n, err := s.Read(buf)
				if err != nil {
					if err != io.EOF {
						fmt.Println("Server: stream read error:", err)
					}
					return
				}

				var ev models.Event
				err = json.Unmarshal(buf[:n], &ev)
				if err != nil {
					fmt.Println("Server: JSON unmarshal error:", err)
					return
				}
				inbound := models.InboundEvent{
					Player: p,
					Ev:     &ev,
				}

				hub.GameEvent <- inbound

				msg := string(buf[:n])
				fmt.Printf("Server: stream %d received '%s'\n", s.StreamID(), msg)

				if _, err := s.Write([]byte("echo: " + msg)); err != nil {
					fmt.Println("Server: write error:", err)
					return
				}
			}
		}(stream, p, s.hub)
	}
}
