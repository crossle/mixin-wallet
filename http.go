package main

import (
	"fmt"
	"net/http"

	"github.com/MixinNetwork/mixin-wallet/middlewares"
	"github.com/dimfeld/httptreemux"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/handlers"
	"github.com/unrolled/render"
)

func StartHTTP() error {
	router := httptreemux.New()
	RegisterHanders(router)
	RegisterRoutes(router)
	handler := middlewares.Authenticate(router)
	handler = middlewares.Constraint(handler)
	handler = middlewares.Context(handler, render.New())
	handler = handlers.ProxyHeaders(handler)

	return gracehttp.Serve(&http.Server{Addr: fmt.Sprintf(":%d", 8001), Handler: handler})
}
