package models

import "github.com/quic-go/quic-go"

type Player struct {
	Id          string      `json:"id"`
	Send        chan []byte `json:"-"`
	QuicStreams QuicStream  `json:"-"`
}

type QuicStream map[string]*quic.Stream
