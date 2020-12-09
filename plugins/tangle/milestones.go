package tangle

import (
	"github.com/massyu/hive.go/workerpool"

	"github.com/massyu/hornet/pkg/model/tangle"
	"github.com/massyu/hornet/plugins/gossip"
)

var (
	processValidMilestoneWorkerCount = 1 // This must not be done in parallel
	processValidMilestoneQueueSize   = 10000
	processValidMilestoneWorkerPool  *workerpool.WorkerPool
)

func processValidMilestone(cachedBndl *tangle.CachedBundle) {
	defer cachedBndl.Release() // bundle -1

	Events.ReceivedNewMilestone.Trigger(cachedBndl) // bundle pass +1

	solidMsIndex := tangle.GetSolidMilestoneIndex()
	bundleMsIndex := cachedBndl.GetBundle().GetMilestoneIndex()

	if tangle.SetLatestMilestoneIndex(bundleMsIndex) {
		Events.LatestMilestoneChanged.Trigger(cachedBndl) // bundle pass +1
		Events.LatestMilestoneIndexChanged.Trigger(bundleMsIndex)
	}
	milestoneSolidifierWorkerPool.TrySubmit(bundleMsIndex, false)
	log.Infof("solidMsIndex is %d", solidMsIndex)
	log.Infof("bundleMsIndex is %d", bundleMsIndex)
	if bundleMsIndex > solidMsIndex {
		log.Infof("Valid milestone detected! Index: %d, Hash: %v", bundleMsIndex, cachedBndl.GetBundle().GetMilestoneHash().Trytes())

		// request trunk and branch
		gossip.RequestMilestoneApprovees(cachedBndl.Retain()) // bundle pass +1
	} else {
		pruningIndex := tangle.GetSnapshotInfo().PruningIndex
		if bundleMsIndex < pruningIndex {
			// this should not happen. we didn't request it and it should be filtered because of timestamp
			log.Warnf("Synced too far back! Index: %d (%v), PruningIndex: %d", bundleMsIndex, cachedBndl.GetBundle().GetMilestoneHash().Trytes(), pruningIndex)
		}
	}
}
