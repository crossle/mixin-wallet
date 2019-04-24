package mixin

import (
	"encoding/hex"

	"github.com/MixinNetwork/go-number"
	"github.com/MixinNetwork/mixin/common"
	"github.com/MixinNetwork/mixin/crypto"
	"github.com/gofrs/uuid"
)

type UTXO struct {
	Asset    string
	Hash     string
	Index    int
	Amount   number.Decimal
	Receiver string
	TraceId  string

	Key  crypto.Key
	Mask crypto.Key
}

func (tx *Transaction) Deposit() *common.DepositData {
	return tx.Inputs[0].Deposit
}

func (tx *Transaction) UTXOs(view crypto.Key) []*UTXO {
	var utxos []*UTXO

	for i, out := range tx.Outputs {
		if out.Type != common.OutputTypeScript {
			continue
		}
		if out.Script != "fffe01" {
			continue
		}
		if len(out.Keys) != 1 {
			continue
		}

		utxo := &UTXO{
			Asset:  tx.Asset,
			Hash:   tx.Hash,
			Index:  i,
			Amount: out.Amount,
			Key:    out.Keys[0],
			Mask:   out.Mask,
		}

		pub := crypto.ViewGhostOutputKey(&utxo.Key, &view, &utxo.Mask, uint64(i))
		utxo.Receiver = pub.String()

		tb, _ := hex.DecodeString(tx.Extra)
		utxo.TraceId = uuid.FromBytesOrNil(tb).String()
		utxos = append(utxos, utxo)
	}

	return utxos
}
