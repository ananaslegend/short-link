package statistic

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	statistic_rows_cap = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "short-link",
		Subsystem: "statistic",
		Name:      "rows_capacity",
		Help:      "how many rows contains in memory of statistic manager. If its value so big better to do flush more often.",
	})
)

func ObserveRowsCap(cap int) {
	statistic_rows_cap.Set(float64(cap))
}
