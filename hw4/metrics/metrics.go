package metrics

import (
	"database/sql"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metr struct {
	Requests *prometheus.CounterVec
	Errors   *prometheus.CounterVec
	Duration *prometheus.SummaryVec
	on       bool
}

const (
	labelQuery    = "request"  // * query text
	labelFunction = "function" // * function to check (Query, QueryRow, Exec)
	labelError    = "error"
)

func New(on bool) *Metr {
	duration := promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "duration_seconds",
			Help:       "Summary of query duration in seconds",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
		},
		[]string{labelQuery, labelFunction, labelError},
	)
	errorsTotal := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors",
		},
		[]string{labelQuery, labelFunction, labelError},
	)

	requestsTotal := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_total",
			Help: "Total number of queries",
		},
		[]string{labelQuery, labelFunction},
	)

	return &Metr{
		Duration: duration,
		Errors:   errorsTotal,
		Requests: requestsTotal,
		on:       on,
	}
}

func (m *Metr) MesurableExec(e func(string, ...interface{}) (sql.Result, error)) func(query string, args ...interface{}) (sql.Result, error) {
	return func(query string, args ...interface{}) (sql.Result, error) {
		t := time.Now()
		if m.on {
			m.Requests.
				WithLabelValues(query, "Exec").
				Inc()
		}
		res, err := e(query, args...)
		if m.on {
			var e string
			if err != nil {
				m.Errors.WithLabelValues(query, "Exec", err.Error()).Inc()
				e = err.Error()
			}
			m.Duration.
				WithLabelValues(query, "Exec", e).
				Observe(time.Since(t).Seconds())
		}
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func (m *Metr) check() {
	db := sql.DB{}
	m.MesurableExec(db.Exec)
}
