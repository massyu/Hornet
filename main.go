package main

import (
	"github.com/iotaledger/hive.go/node"

	"github.com/massyu/hornet/pkg/config"
	"github.com/massyu/hornet/pkg/toolset"
	"github.com/massyu/hornet/plugins/autopeering"
	"github.com/massyu/hornet/plugins/cli"
	"github.com/massyu/hornet/plugins/coordinator"
	"github.com/massyu/hornet/plugins/dashboard"
	"github.com/massyu/hornet/plugins/database"
	"github.com/massyu/hornet/plugins/gossip"
	"github.com/massyu/hornet/plugins/gracefulshutdown"
	"github.com/massyu/hornet/plugins/metrics"
	"github.com/massyu/hornet/plugins/mqtt"
	"github.com/massyu/hornet/plugins/peering"
	"github.com/massyu/hornet/plugins/pow"
	"github.com/massyu/hornet/plugins/profiling"
	"github.com/massyu/hornet/plugins/prometheus"
	"github.com/massyu/hornet/plugins/snapshot"
	"github.com/massyu/hornet/plugins/spammer"
	"github.com/massyu/hornet/plugins/tangle"
	"github.com/massyu/hornet/plugins/urts"
	"github.com/massyu/hornet/plugins/warpsync"
	"github.com/massyu/hornet/plugins/webapi"
	"github.com/massyu/hornet/plugins/zmq"
)

func main() {
	cli.HideConfigFlags()
	cli.ParseFlags()
	cli.PrintVersion()
	cli.ParseConfig()
	toolset.HandleTools()
	cli.PrintConfig()

	plugins := []*node.Plugin{
		cli.PLUGIN,
		gracefulshutdown.PLUGIN,
		profiling.PLUGIN,
		database.PLUGIN,
		autopeering.PLUGIN,
		webapi.PLUGIN,
	}

	if !config.NodeConfig.GetBool(config.CfgNetAutopeeringRunAsEntryNode) {
		plugins = append(plugins, []*node.Plugin{
			pow.PLUGIN,
			gossip.PLUGIN,
			tangle.PLUGIN,
			peering.PLUGIN,
			warpsync.PLUGIN,
			urts.PLUGIN,
			metrics.PLUGIN,
			snapshot.PLUGIN,
			dashboard.PLUGIN,
			zmq.PLUGIN,
			mqtt.PLUGIN,
			spammer.PLUGIN,
			coordinator.PLUGIN,
			prometheus.PLUGIN,
		}...)
	}

	node.Run(node.Plugins(plugins...))
}
