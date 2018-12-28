package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	msgsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ingester_msgs_processed",
		Help: "The total number of processed messages",
	})
	pointsCommitted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ingester_points_committed",
		Help: "The total number of committed values",
	})
	activeSubscriptions = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ingester_active_subscriptions",
		Help: "# of active WAVEMQ subscriptions",
	})
)
