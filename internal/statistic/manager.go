package statistic

import (
	"context"
	"errors"
	"fmt"
	"github.com/ananaslegend/short-link/pkg/logs"
	"log/slog"
	"sync"
	"time"
)

// Writer should insert stats to db.
type Writer interface {
	InsertRows(Rows) error
}

type Manager struct {
	log *slog.Logger

	mx      sync.RWMutex
	rows    Rows
	rowsCap int

	flushTime time.Duration
	writer    Writer

	shutdownCh    chan struct{}
	shutdownErrCh chan error
}

func NewManager(flushTime time.Duration, rowsCap int, repository Writer, log *slog.Logger) *Manager {
	return &Manager{
		rows:          newRows(rowsCap),
		rowsCap:       rowsCap,
		flushTime:     flushTime,
		shutdownCh:    make(chan struct{}),
		shutdownErrCh: make(chan error),
		writer:        repository,
		log:           log,
	}
}

func (m *Manager) Append(dimension Dimension, metric Metric) {
	m.mx.Lock()
	defer m.mx.Unlock()

	current := m.rows[dimension]
	current = current.append(metric)

	m.rows[dimension] = current
}

func (m *Manager) AppendRow(row *Row) {
	m.Append(row.Dimension, row.Metric)
}

func (m *Manager) insert() error {
	const op = "statistic.manager.insert"
	rowsToInsert := m.withdrawRows()
	if len(rowsToInsert) == 0 {
		return ErrNoSatatToInsert
	}

	if err := m.writer.InsertRows(rowsToInsert); err != nil {
		for dimension, metric := range rowsToInsert {
			m.Append(dimension, metric)
		}
		return fmt.Errorf("%v: %v", op, err)
	}
	return nil
}

func (m *Manager) withdrawRows() Rows {
	m.mx.Lock()
	defer m.mx.Unlock()

	rows := m.rows
	m.rows = newRows(m.rowsCap)

	return rows
}

func (m *Manager) Close(ctx context.Context) error {
	m.shutdownCh <- struct{}{}
	defer close(m.shutdownCh)

	err := <-m.shutdownErrCh
	defer close(m.shutdownErrCh)

	return err
}

func (m *Manager) loop() {
	for {
		select {
		case <-time.After(m.flushTime):
			if err := m.insert(); err != nil {
				switch {
				case errors.Is(err, ErrNoSatatToInsert):
					m.log.Warn(ErrNoSatatToInsert.Error())
				default:
					m.log.Error("cant insert stat rows in db", logs.Err(err))
				}
			}
		case <-m.shutdownCh:
			m.shutdownErrCh <- m.insert()
			return
		}
	}
}

func (m *Manager) Run() {
	go m.loop()
}
