package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MixinNetwork/mixin-wallet/mixin"
	"github.com/MixinNetwork/mixin-wallet/models"
)

type SnapService struct{}

func (service *SnapService) Run(ctx context.Context) error {
	service.loopSnap(ctx)
	return nil
}

func (service *SnapService) loopSnap(ctx context.Context) {
	rpc := mixin.NewMixinNetwork(node)
	for {
		checkpoint, err := readSnapshotCheckPoint(ctx, "key")
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
			if err := models.CreateSnapshot(ctx, s.Hash, int64(s.Topology), int64(s.Timestamp), tx.Hash); err != nil {
				log.Println(err)
				break
			}
			c := fmt.Sprintf("%d", checkpoint)
			fmt.Println(c)
			err = models.WriteProperty(ctx, "key", c)
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
