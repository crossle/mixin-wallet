package main

import (
	"fmt"
	"net/http"

	"github.com/crossle/mixin-wallet/durable"
	"github.com/crossle/mixin-wallet/middlewares"
	"github.com/dimfeld/httptreemux"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/handlers"
	"github.com/unrolled/render"
)

func StartHTTP(db *durable.Database) error {
	router := httptreemux.New()
	RegisterHanders(router)
	RegisterRoutes(router)
	handler := middlewares.Authenticate(router)
	handler = middlewares.Constraint(handler)
	handler = middlewares.Context(handler, db, render.New())
	handler = handlers.ProxyHeaders(handler)

	return gracehttp.Serve(&http.Server{Addr: fmt.Sprintf(":%d", 8001), Handler: handler})
}
