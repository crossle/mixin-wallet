package models

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/MixinNetwork/mixin-wallet/durable"
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
CREATE INDEX IF NOT EXISTS index_snapshots_transaction_hash ON snapshots(transaction_hash);
`

type Snapshot struct {
	Hash            string `json:"hash"`
	Topology        int64  `json:"topology"`
	Timestamp       int64  `json:"timestamp"`
	TransactionHash string `json:"transaction_hash"`
}

func CreateSnapshot(ctx context.Context, hash string, topology, timestamp int64, transactionHash string) error {
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, txn *sql.Tx) error {
		query := "SELECT hash, transaction_hash, topology, timestamp FROM snapshots WHERE hash=$1"
		row := txn.QueryRowContext(ctx, query, hash)
		_, err := snapshotFromRow(row)
		switch {
		case err == sql.ErrNoRows:
			query := "INSERT INTO snapshots (hash, transaction_hash, topology, timestamp) VALUES ($1, $2, $3, $4)"
			_, err := txn.ExecContext(ctx, query, hash, transactionHash, topology, time.Unix(0, timestamp))
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

func QuerySnapshotByHash(ctx context.Context, hash string) (*Snapshot, error) {
	query := "SELECT hash, transaction_hash, topology, timestamp FROM snapshots WHERE hash=$1"
	row := session.Database(ctx).QueryRowContext(ctx, query, hash)
	snapshot, err := snapshotFromRow(row)
	if err == sql.ErrNoRows {
		return &Snapshot{}, nil
	} else if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return snapshot, nil
}

func QuerySnapshotByTransactionHash(ctx context.Context, transactionHash string) (*Snapshot, error) {
	query := "SELECT hash, transaction_hash, topology, timestamp FROM snapshots WHERE transaction_hash=$1"
	row := session.Database(ctx).QueryRowContext(ctx, query, transactionHash)
	snapshot, err := snapshotFromRow(row)
	if err == sql.ErrNoRows {
		return &Snapshot{}, nil
	} else if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return snapshot, nil
}

func snapshotFromRow(row durable.Row) (*Snapshot, error) {
	var s Snapshot
	err := row.Scan(&s.Hash, &s.TransactionHash, &s.Topology, &s.Timestamp)
	return &s, err
}
