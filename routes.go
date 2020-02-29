package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/crossle/mixin-wallet/mixin"
	"github.com/crossle/mixin-wallet/models"
	"github.com/crossle/mixin-wallet/views"
	"github.com/dimfeld/httptreemux"
)

var node string

func init() {
	node = "http://mixin-node-02.b1.run:8239"
}

func RegisterRoutes(router *httptreemux.TreeMux) {
	router.GET("/createaddress", createAddress)
	router.GET("/getinfo", getNodeInfo)
	router.GET("/listallnodes", getAllNodes)
	router.GET("/height", getHeight)
	router.GET("/snapshots", getSnapshots)
	router.GET("/snapshots/:id", getSnapshot)
	router.GET("/transactions/:id/utxo", getTransactionUTXO)
	router.GET("/transactions/:id", getTransaction)
	router.GET("/transactions/:id/snapshot", getTransactionSnapshot)
	router.POST("/transactions", postRaw)
	router.GET("/account/:id", getAccount)
}

func getAccount(w http.ResponseWriter, r *http.Request, params map[string]string) {

}

func createAddress(w http.ResponseWriter, r *http.Request, params map[string]string) {
	key, err := mixin.LocalGenerateKey()
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderDataResponse(w, r, map[string]interface{}{"address": key.Address, "view": key.ViewKey, "spend": key.SpendKey})
}

func getNodeInfo(w http.ResponseWriter, r *http.Request, params map[string]string) {
	n := r.URL.Query().Get("node")
	if n == "" {
		n = node
	}
	rpc := mixin.NewMixinNetwork(n)
	nodeInfo, err := rpc.GetInfo()
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderDataResponse(w, r, nodeInfo)
}

func getAllNodes(w http.ResponseWriter, r *http.Request, params map[string]string) {
	n := r.URL.Query().Get("node")
	if n == "" {
		n = node
	}
	rpc := mixin.NewMixinNetwork(n)
	nodes, err := rpc.ListAllNodes()
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderDataResponse(w, r, nodes)
}

func getHeight(w http.ResponseWriter, r *http.Request, params map[string]string) {
	rpc := mixin.NewMixinNetwork(node)
	nodeInfo, err := rpc.GetInfo()
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderDataResponse(w, r, map[string]interface{}{"height": nodeInfo.Graph.Topology})
}

func getSnapshots(w http.ResponseWriter, r *http.Request, params map[string]string) {
	rpc := mixin.NewMixinNetwork(node)
	s := r.URL.Query().Get("offset")
	c := r.URL.Query().Get("limit")
	var since, count uint64
	var err error
	if s != "" {
		since, err = strconv.ParseUint(s, 10, 64)
		if err != nil {
			views.RenderErrorResponse(w, r, err)
			return
		}
	}
	if c != "" {
		count, err = strconv.ParseUint(c, 10, 64)
		if err != nil {
			views.RenderErrorResponse(w, r, err)
			return
		}
	}
	snapshots, err := rpc.ListSnapshotsSince(since, count)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderDataResponse(w, r, snapshots)
}

func getSnapshot(w http.ResponseWriter, r *http.Request, params map[string]string) {
	rpc := mixin.NewMixinNetwork(node)
	id := params["id"]
	since, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		snapshot, err := models.QuerySnapshotByHash(r.Context(), id)
		if err != nil {
			views.RenderErrorResponse(w, r, err)
			return
		}
		since = uint64(snapshot.Topology)
	}
	snapshots, err := rpc.ListSnapshotsSince(since, 1)
	if err != nil || len(snapshots) != 1 {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderDataResponse(w, r, snapshots[0])
}

func getTransaction(w http.ResponseWriter, r *http.Request, params map[string]string) {
	n := r.URL.Query().Get("node")
	if n == "" {
		n = node
	}
	rpc := mixin.NewMixinNetwork(n)
	transaction, err := rpc.GetTransaction(params["id"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderDataResponse(w, r, transaction)
}

func getTransactionSnapshot(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id := params["id"]
	snapshot, err := models.QuerySnapshotByTransactionHash(r.Context(), id)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderDataResponse(w, r, snapshot)
}

func getTransactionUTXO(w http.ResponseWriter, r *http.Request, params map[string]string) {
	rpc := mixin.NewMixinNetwork(node)
	viewKey := r.URL.Query().Get("view")
	utxos, err := rpc.GetTransactionUTXO(params["id"], viewKey)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderDataResponse(w, r, utxos)
}

func postRaw(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body struct {
		Raw string `json:"raw"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	rpc := mixin.NewMixinNetwork(node)
	txId, err := rpc.SendRawTransaction(body.Raw)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderDataResponse(w, r, map[string]interface{}{"transaction_hash": txId})
}
