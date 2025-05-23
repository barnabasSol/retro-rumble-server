package server

import "github.com/barnabasSol/retro-rumble/internals/event"

type Handlers map[string]func(ev event.GameEvent, client Client)
