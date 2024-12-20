package closer

import (
	"strings"
	"sync"
)

type messages struct {
	mx  sync.Mutex
	msg []string
}

func newMessages(srvCount int) *messages {
	return &messages{
		mx:  sync.Mutex{},
		msg: make([]string, 0, srvCount),
	}
}

func (m *messages) add(msg string) {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.msg = append(m.msg, msg)
}

func (m *messages) len() int {
	m.mx.Lock()
	defer m.mx.Unlock()

	return len(m.msg)
}

func (m *messages) String() string {
	m.mx.Lock()
	defer m.mx.Unlock()

	return strings.Join(m.msg, "\n")
}
