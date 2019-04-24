package main

import (
	"net/http"

	"github.com/MixinNetwork/mixin-wallet/session"
	"github.com/MixinNetwork/mixin-wallet/views"
	"github.com/dimfeld/httptreemux"
)

func RegisterHanders(router *httptreemux.TreeMux) {
	router.MethodNotAllowedHandler = func(w http.ResponseWriter, r *http.Request, _ map[string]httptreemux.HandlerFunc) {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	}
	router.NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	}
}
