package statistic

import (
	"time"
)

type Dimension struct {
	Timestamp int64
	Link      string // TODO: Order instead
}

func NewDimension() Dimension {
	timestamp := time.Now().Unix()
	timestamp -= timestamp % 60

	return Dimension{
		Timestamp: timestamp,
	}
}

type Metric struct {
	Redirect int
}

type Row struct {
	Dimension
	Metric
}

func NewRow() *Row {
	return &Row{
		Dimension: NewDimension(),
		Metric:    Metric{},
	}
}

type Rows map[Dimension]Metric

func newRows(rowsCap int) Rows {
	return make(Rows, rowsCap)
}

func (m Metric) append(new Metric) Metric {
	m.Redirect += new.Redirect

	return m
}
