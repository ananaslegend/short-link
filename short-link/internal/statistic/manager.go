package statistic

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/ananaslegend/short-link/pkg/cslog"
)

// Writer should insert stats to db.
type Writer interface {
	InsertRows(context.Context, Rows) error
}

type Manager struct {
	log *slog.Logger

	mx      sync.RWMutex
	rows    Rows
	rowsCap int

	flushTime time.Duration
	writer    Writer

	shutdownCh    chan context.Context
	shutdownErrCh chan error
}

func NewManager(flushTime time.Duration, rowsCap int, repository Writer, log *slog.Logger) *Manager {
	return &Manager{
		rows:          newRows(rowsCap),
		rowsCap:       rowsCap,
		flushTime:     flushTime,
		shutdownCh:    make(chan context.Context),
		shutdownErrCh: make(chan error),
		writer:        repository,
		log:           log,
	}
}

func (m *Manager) Close(ctx context.Context) error {
	m.shutdownCh <- ctx
	defer close(m.shutdownCh)

	err := <-m.shutdownErrCh
	defer close(m.shutdownErrCh)

	return err
}

func (m *Manager) Append(dimension Dimension, metric Metric) {
	m.mx.Lock()
	defer m.mx.Unlock()

	current := m.rows[dimension]
	current = current.append(metric)

	m.rows[dimension] = current
	ObserveRowsCap(len(m.rows))
}

func (m *Manager) AppendRow(row *Row) {
	m.Append(row.Dimension, row.Metric)
}

func (m *Manager) insert(ctx context.Context) error {
	const op = "statistic.manager.insert"
	rowsToInsert := m.withdrawRows()
	if len(rowsToInsert) == 0 {
		return ErrNoStatToInsert
	}

	if err := m.writer.InsertRows(ctx, rowsToInsert); err != nil {
		for dimension, metric := range rowsToInsert {
			m.Append(dimension, metric)
		}

		return fmt.Errorf("%v: %v", op, err)
	}

	ObserveRowsCap(len(m.rows))

	return nil
}

func (m *Manager) withdrawRows() Rows {
	m.mx.Lock()
	defer m.mx.Unlock()

	rows := m.rows
	m.rows = newRows(m.rowsCap)

	return rows
}

func (m *Manager) loop() {
	for {
		select {
		case <-time.After(m.flushTime):
			if err := m.insert(context.Background()); err != nil {
				switch {
				case errors.Is(err, ErrNoStatToInsert):
					m.log.Info(ErrNoStatToInsert.Error())
				default:
					m.log.Error("cant insert stat rows in db", cslog.Error(err))
				}
			} else {
				m.log.Debug("statistic flushed")
			}
		case ctx := <-m.shutdownCh:
			m.shutdownErrCh <- m.insert(ctx)
			return
		}
	}
}

func (m *Manager) Run() {
	m.loop()
}

func (m *Manager) FlushTime() time.Duration {
	return m.flushTime
}
