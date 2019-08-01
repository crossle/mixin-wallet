package services

import (
	"context"
	"fmt"

	"github.com/crossle/mixin-wallet/durable"
	"github.com/crossle/mixin-wallet/session"
)

type Hub struct {
	context  context.Context
	services map[string]Service
}

func NewHub(db *durable.Database) *Hub {
	hub := &Hub{services: make(map[string]Service)}
	hub.context = session.WithDatabase(context.Background(), db)
	hub.registerServices()
	return hub
}

func (hub *Hub) StartService(name string) error {
	service := hub.services[name]
	if service == nil {
		return fmt.Errorf("no service found: %s", name)
	}

	return service.Run(hub.context)
}

func (hub *Hub) registerServices() {
	hub.services["scan"] = &ScanService{}
}
