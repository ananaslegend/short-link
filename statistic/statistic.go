package statistic

import (
	"context"
	"time"
)

type Dimension struct {
	Timestamp int64
	Link      string
	Alias     string
}

func NewDimension() Dimension {
	timestamp := time.Now().Unix()

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

func (r *Row) IsEmpty() bool {
	if r.Metric.Redirect == 0 {
		return true
	}

	return false
}

type Rows map[Dimension]Metric

func newRows(rowsCap int) Rows {
	return make(Rows, rowsCap)
}

func (m Metric) append(new Metric) Metric {
	m.Redirect += new.Redirect

	return m
}

type statRowCtxKey struct{}

type StatManager interface {
	AppendRow(row *Row)
	FlushTime() time.Duration
}

func (r *Row) SetToCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, statRowCtxKey{}, r)
}

func SetToCtx(ctx context.Context, r *Row) context.Context {
	return r.SetToCtx(ctx)
}

func GetFromCtx(ctx context.Context) (*Row, bool) {
	r, ok := ctx.Value(statRowCtxKey{}).(*Row)

	return r, ok
}
