package models

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/MixinNetwork/mixin-wallet/mixin"
	"github.com/MixinNetwork/mixin-wallet/session"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
)

const utxos_DDL = `
CREATE TABLE IF NOT EXISTS utxos (
	utxo_id VARCHAR(36) PRIMARY KEY,
	amount VARCHAR(128) NOT NULL,
	asset_id VARCHAR(64) NOT NULL,
	output_index INTEGER NOT NULL,
	output_script VARCHAR(36) NOT NULL,
	output_type  INTEGER NOT NULL,
	state VARCHAR(36) NOT NULL,
	spent_by VARCHAR(64),
	keys text[] NOT NULL,
	mask VARCHAR(128),
	extra VARCHAR(128),
	transaction_hash VARCHAR(64) NOT NULL,
	created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
`

type UTXO struct {
	UTXOId          string
	TransactionHash string
	AssetId         string
	Amount          string
	OutputIndex     int
	OutputScript    string
	OutputType      uint8
	State           string
	Keys            []string
	Mask            string
	Extra           string
	SpentBy         string
}

const (
	UTXOStateUnspent = "unspent"
	UTXOStateSpent   = "spent"
)

func queryForRow(utxoId string) {

}

func CreateOrUpdateUTXOs(ctx context.Context, tx *mixin.Transaction) error {
	for _, input := range tx.Inputs {
		if input.Hash == "" {
			continue
		}
		utxoId := uniqueMixinUTXOId(input.Hash, input.Index)
		err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, txn *sql.Tx) error {
			var u UTXO
			query := "SELECT utxo_id FROM utxos WHERE utxo_id=$1"
			if err := txn.QueryRowContext(ctx, query, utxoId).Scan(&u.UTXOId); err != nil {
				return err
			}
			if u.OutputIndex != input.Index {
				return fmt.Errorf("erorr transaction index: %s %d", input.Hash, input.Index)
			}
			_, err := txn.ExecContext(ctx, "UPDATE utxos SET spent_by = $1, state = $2 WHERE utxo_id = $3", tx.Hash, UTXOStateSpent, utxoId)
			return err
		})
		if err != nil {
			return session.TransactionError(ctx, err)
		}
	}
	var utxos []*UTXO
	for index, output := range tx.Outputs {
		utxoId := uniqueMixinUTXOId(tx.Hash, index)
		utxo := UTXO{
			UTXOId:          utxoId,
			TransactionHash: tx.Hash,
			AssetId:         tx.Asset,
			Amount:          output.Amount.String(),
			OutputIndex:     index,
			OutputScript:    output.Script,
			OutputType:      output.Type,
			State:           UTXOStateUnspent,
			Mask:            output.Mask.String(),
			Extra:           tx.Extra,
		}
		var keys []string
		for _, k := range output.Keys {
			e := base64.StdEncoding.EncodeToString(k[:])
			keys = append(keys, e)
		}
		utxo.Keys = keys
		utxos = append(utxos, &utxo)
	}

	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, txn *sql.Tx) error {
		stmt, err := txn.Prepare(pq.CopyIn("utxos", "utxo_id", "transaction_hash", "asset_id", "amount", "output_index", "output_script", "output_type", "state", "mask", "extra", "keys"))
		if err != nil {
			return err
		}

		for _, utxo := range utxos {
			_, err = stmt.Exec(utxo.UTXOId, utxo.TransactionHash, utxo.AssetId, utxo.Amount, utxo.OutputIndex, utxo.OutputScript, utxo.OutputType, utxo.State, utxo.Mask, utxo.Extra, pq.Array(utxo.Keys))
			if err != nil {
				return err
			}
		}
		_, err = stmt.Exec()
		if err != nil {
			return err
		}
		err = stmt.Close()
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}

func uniqueMixinUTXOId(hash string, index int) string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s:%d", hash, index)))
	s := h.Sum(nil)
	s[6] = (s[6] & 0x0f) | 0x30
	s[8] = (s[8] & 0x3f) | 0x80
	sid, err := uuid.FromBytes(s)
	if err != nil {
		panic(err)
	}
	return sid.String()
}
