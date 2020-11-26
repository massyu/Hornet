package metrics

import (
	"time"

	"github.com/massyu/hive.go/daemon"
	"github.com/massyu/hive.go/node"
	"github.com/massyu/hive.go/timeutil"

	"github.com/massyu/hornet/pkg/shutdown"
)

var PLUGIN = node.NewPlugin("Metrics", node.Enabled, configure, run)

func configure(_ *node.Plugin) {
	// nothing
}

func run(_ *node.Plugin) {
	// create a background worker that "measures" the TPS value every second
	daemon.BackgroundWorker("Metrics TPS Updater", func(shutdownSignal <-chan struct{}) {
		timeutil.Ticker(measureTPS, 1*time.Second, shutdownSignal)
	}, shutdown.PriorityMetricsUpdater)
}
