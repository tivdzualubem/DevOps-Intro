package main

import (
	"fmt"
	"io"
	"strconv"
	"sync/atomic"
	"time"
)

const requestDurationMetricName = "quicknotes_http_request_duration_seconds"

var requestDurationBuckets = []float64{
	0.00005,
	0.0001,
	0.00025,
	0.0005,
	0.001,
	0.0025,
	0.005,
	0.01,
	0.025,
	0.05,
	0.1,
	0.25,
	0.5,
	1,
	2,
	5,
}

type requestDurationMetrics struct {
	bucketCounts   []atomic.Uint64
	count          atomic.Uint64
	sumNanoseconds atomic.Uint64
}

func newRequestDurationMetrics() *requestDurationMetrics {
	return &requestDurationMetrics{
		bucketCounts: make(
			[]atomic.Uint64,
			len(requestDurationBuckets),
		),
	}
}

func (m *requestDurationMetrics) observe(duration time.Duration) {
	seconds := duration.Seconds()

	m.count.Add(1)
	m.sumNanoseconds.Add(uint64(duration.Nanoseconds()))

	for index, upperBound := range requestDurationBuckets {
		if seconds <= upperBound {
			m.bucketCounts[index].Add(1)
		}
	}
}

func (m *requestDurationMetrics) writePrometheus(writer io.Writer) {
	_, _ = fmt.Fprintf(
		writer,
		"# HELP %s HTTP request duration in seconds.\n",
		requestDurationMetricName,
	)
	_, _ = fmt.Fprintf(
		writer,
		"# TYPE %s histogram\n",
		requestDurationMetricName,
	)

	for index, upperBound := range requestDurationBuckets {
		label := strconv.FormatFloat(
			upperBound,
			'g',
			-1,
			64,
		)

		_, _ = fmt.Fprintf(
			writer,
			"%s_bucket{le=\"%s\"} %d\n",
			requestDurationMetricName,
			label,
			m.bucketCounts[index].Load(),
		)
	}

	count := m.count.Load()
	sumSeconds :=
		float64(m.sumNanoseconds.Load()) /
			float64(time.Second)

	_, _ = fmt.Fprintf(
		writer,
		"%s_bucket{le=\"+Inf\"} %d\n",
		requestDurationMetricName,
		count,
	)
	_, _ = fmt.Fprintf(
		writer,
		"%s_sum %s\n",
		requestDurationMetricName,
		strconv.FormatFloat(sumSeconds, 'g', -1, 64),
	)
	_, _ = fmt.Fprintf(
		writer,
		"%s_count %d\n",
		requestDurationMetricName,
		count,
	)
}
