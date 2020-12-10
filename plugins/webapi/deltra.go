package webapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/massyu/hornet/pkg/config"
	"github.com/massyu/hornet/pkg/model/tangle"
	"github.com/massyu/hornet/plugins/cli"
)

func init() {
	addEndpoint("deleteTransaction", deleteTransaction, implementedAPIcalls)
	addEndpoint("deleteAPIConfiguration", deleteAPIConfiguration, implementedAPIcalls)
}

func deleteTransaction(_ interface{}, c *gin.Context, _ <-chan struct{}) {
	// Basic info data
	result := DeleteTransactionReturn{
		AppName:    cli.AppName,
		AppVersion: cli.AppVersion,
	}

	// Return node info
	c.JSON(http.StatusOK, result)
}

func deleteAPIConfiguration(_ interface{}, c *gin.Context, _ <-chan struct{}) {

	result := DeleteTransactionConfigurationReturn{
		MaxFindTransactions: config.NodeConfig.GetInt(config.CfgWebAPILimitsMaxFindTransactions),
		MaxRequestsList:     config.NodeConfig.GetInt(config.CfgWebAPILimitsMaxRequestsList),
		MaxGetTrytes:        config.NodeConfig.GetInt(config.CfgWebAPILimitsMaxGetTrytes),
		MaxBodyLength:       config.NodeConfig.GetInt(config.CfgWebAPILimitsMaxBodyLengthBytes),
	}

	// Milestone start index
	snapshotInfo := tangle.GetSnapshotInfo()
	if snapshotInfo != nil {
		result.MilestoneStartIndex = snapshotInfo.PruningIndex
	}

	c.JSON(http.StatusOK, result)
}
