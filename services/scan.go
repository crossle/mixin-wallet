package services

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/MixinNetwork/mixin-wallet/mixin"
	"github.com/MixinNetwork/mixin-wallet/models"
)

type ScanService struct{}

var node string

func init() {
	node = "http://35.188.235.212:8239"
}

func (service *ScanService) Run(ctx context.Context) error {
	service.loopSnapshots(ctx)
	return nil
}

func (service *ScanService) loopSnapshots(ctx context.Context) {
	rpc := mixin.NewMixinNetwork(node)
	count := uint64(100)
	for {
		checkpoint, err := readSnapshotCheckPoint(ctx)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		snapshots, err := rpc.ListSnapshotsSince(checkpoint, count)
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}
		for _, s := range snapshots {
			checkpoint = s.Topology
			tx := &s.Transaction
			if err := models.CreateOrUpdateUTXOs(ctx, tx, s.Timestamp); err != nil {
				log.Println(err)
				break
			}
			c := fmt.Sprintf("%d", checkpoint)
			fmt.Println(c)
			err = models.WriteProperty(ctx, models.MainnetSnapshotsCheckpoint, c)
			if err != nil {
				log.Println(err)
				break
			}
		}
		if uint64(len(snapshots)) < count {
			time.Sleep(1 * time.Second)
		}
	}
}

func readSnapshotCheckPoint(ctx context.Context) (uint64, error) {
	since, err := models.ReadProperty(ctx, models.MainnetSnapshotsCheckpoint)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	if since == "" {
		since = "0"
	}
	h, err := strconv.ParseUint(since, 10, 64)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return h, nil
}
