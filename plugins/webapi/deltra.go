package webapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/massyu/hornet/pkg/config"
	"github.com/massyu/hornet/pkg/model/tangle"
	"github.com/massyu/hornet/plugins/cli"
	"github.com/massyu/hornet/plugins/coordinator"
	"github.com/mitchellh/mapstructure"
)

func init() {
	addEndpoint("deleteTransaction", deleteTransaction, implementedAPIcalls)
	addEndpoint("deleteAPIConfiguration", deleteAPIConfiguration, implementedAPIcalls)
}

func deleteTransaction(i interface{}, c *gin.Context, _ <-chan struct{}) {
	e := ErrorReturn{}
	query := &DeleteTransaction{}
	if err := mapstructure.Decode(i, query); err != nil {
		e.Error = fmt.Sprintf("%v: %v", ErrInternalError, err)
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	log.Info(query.Bundle)
	log.Info(query.Command)
	// Basic info data
	result := DeleteTransactionReturn{
		AppName:       cli.AppName,
		BundleAddress: query.Bundle,
	}

	coordinator.SetCancelSignal(query.Bundle)

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
