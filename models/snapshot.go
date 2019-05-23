package models

import (
	"context"
	"database/sql"
	"log"

	"github.com/MixinNetwork/mixin-wallet/session"
)

const snapshots_DDL = `
CREATE TABLE IF NOT EXISTS snapshots (
	hash VARCHAR(64) PRIMARY KEY,
	transaction_hash VARCHAR(64) NOT NULL,
	topology INTEGER NOT NULL,
	timestamp        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS index_snapshots_topology ON snapshots(topology);
`

type Snapshot struct {
	Hash             string
	topology         int64
	timestamp        int64
	transaction_hash string
}

func CreateSnapshot(ctx context.Context, hash string, topology, timestamp int64, transactionHash string) error {
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, txn *sql.Tx) error {
		var s Snapshot
		query := "SELECT hash FROM snapshots WHERE hash=$1"
		err := txn.QueryRowContext(ctx, query, hash).Scan(&s.Hash)
		switch {
		case err == sql.ErrNoRows:
			query := "INSERT INTO snapshots (hash, transaction_hash, topology, timestamp) VALUES ($1, $2, $3, $4)"
			_, err := txn.ExecContext(ctx, query, hash, transactionHash, topology, timestamp)
			if err != nil {
				return err
			}
		case err != nil:
			log.Fatalf("query error: %v\n", err)
			return err
		default:
		}
		return nil
	})
	return err
}
