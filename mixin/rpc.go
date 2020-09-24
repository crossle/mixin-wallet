package mixin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MixinNetwork/go-number"
	"github.com/MixinNetwork/mixin/common"
	"github.com/MixinNetwork/mixin/crypto"
)

type MixinNetwork struct {
	httpClient *http.Client
	node       string
}

type Transaction struct {
	Version uint8  `json:"version"`
	Asset   string `json:"asset"`
	Inputs  []struct {
		Hash    string              `json:"hash"`
		Index   int                 `json:"index"`
		Genesis string              `json:"genesis"`
		Deposit *common.DepositData `json:"deposit"`
		Mint    *common.MintData    `json:"mint"`
	} `json:"inputs"`
	Outputs []struct {
		Type   uint8          `json:"type"`
		Amount number.Decimal `json:"amount"`
		Keys   []crypto.Key   `json:"keys"`
		Script string         `json:"script"`
		Mask   crypto.Key     `json:"mask"`
	} `json:"outputs"`
	Extra string `json:"extra"`
	Hash  string `json:"hash"`
}

type SnapshotWithTransaction struct {
	Hash        string      `json:"hash"`
	Timestamp   uint64      `json:"timestamp"`
	Topology    uint64      `json:"topology"`
	Transaction Transaction `json:"transaction"`
}

type Node struct {
	Id          string `json:"id`
	Signer      string `json:"signer"`
	Payee       string `json:"payee"`
	Transaction string `json:"transaction"`
	Timestamp   int64  `json:"timestamp"`
	State       string `json:"state"`
}

type NodeInfo struct {
	Network string `json:"network"`
	Node    string `json:"node"`
	Version string `json:"version"`
	Uptime  string `json:"uptime"`
	Queue   struct {
		Caches       int64 `json:"caches"`
		Finals       int64 `json:"finals"`
		Transactions int64 `json:"transactions"`
	} `json:"queue"`
	Graph struct {
		Topology  int64 `json:"topology"`
		Consensus []struct {
			Node      string `json:"node"`
			Payee     string `json:"payee"`
			Signer    string `json:"signer"`
			State     string `json:"state"`
			Timestamp int64  `json:"timestamp"`
		} `json:"consensus"`
	} `json:"graph"`
}

func NewMixinNetwork(node string) *MixinNetwork {
	return &MixinNetwork{
		httpClient: &http.Client{Timeout: 20 * time.Second},
		node:       node,
	}
}

func (m *MixinNetwork) SendRawTransaction(raw string) (string, error) {
	body, err := m.callRPC("sendrawtransaction", []interface{}{raw})
	if err != nil {
		return "", err
	}
	var tx Transaction
	err = json.Unmarshal(body, &tx)
	return tx.Hash, err
}

func (m *MixinNetwork) GetTransactionUTXO(hash string, viewKey string) ([]*UTXO, error) {
	if viewKey == "" {
		return nil, fmt.Errorf("No view key for this transaction %s", hash)
	}
	tx, err := m.GetTransaction(hash)
	if err != nil {
		return nil, err
	}
	key, err := ParseKeyFromHex(viewKey)
	if err != nil {
		return nil, err
	}
	utxos := tx.UTXOs(key)
	return utxos, err
}

func (m *MixinNetwork) GetTransaction(hash string) (*Transaction, error) {
	body, err := m.callRPC("gettransaction", []interface{}{hash})
	if err != nil {
		return nil, err
	}
	var tx Transaction
	err = json.Unmarshal(body, &tx)
	if err != nil || tx.Hash == "" {
		return nil, err
	}
	return &tx, err
}

func (m *MixinNetwork) ListSnapshotsSince(since, count uint64) ([]*SnapshotWithTransaction, error) {
	if count == 0 {
		count = 10
	}
	body, err := m.callRPC("listsnapshots", []interface{}{since, count, 0, 1})
	if err != nil {
		return nil, err
	}
	var snapshots []*SnapshotWithTransaction
	err = json.Unmarshal(body, &snapshots)
	return snapshots, err
}

func (m *MixinNetwork) GetInfo() (*NodeInfo, error) {
	body, err := m.callRPC("getinfo", []interface{}{})
	if err != nil {
		return nil, err
	}
	var nodeInfo NodeInfo
	err = json.Unmarshal(body, &nodeInfo)
	if err != nil {
		return nil, err
	}
	return &nodeInfo, nil
}

func (m *MixinNetwork) ListAllNodes() ([]*Node, error) {
	body, err := m.callRPC("listallnodes", []interface{}{})
	if err != nil {
		return nil, err
	}
	var nodes []*Node
	err = json.Unmarshal(body, &nodes)
	return nodes, err
}

func (m *MixinNetwork) callRPC(method string, params []interface{}) ([]byte, error) {
	body, err := json.Marshal(map[string]interface{}{
		"method": method,
		"params": params,
	})
	if err != nil {
		return nil, err
	}
	endpoint := "http://" + m.node
	if strings.HasPrefix(m.node, "http") {
		endpoint = m.node
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data  interface{} `json:"data"`
		Error interface{} `json:"error"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, fmt.Errorf("ERROR %s", result.Error)
	}

	return json.Marshal(result.Data)
}
