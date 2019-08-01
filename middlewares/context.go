package middlewares

import (
	"net/http"

	"github.com/crossle/mixin-wallet/durable"
	"github.com/crossle/mixin-wallet/session"
	"github.com/unrolled/render"
)

func Context(handler http.Handler, db *durable.Database, render *render.Render) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := session.WithRequest(r.Context(), r)
		ctx = session.WithDatabase(ctx, db)
		ctx = session.WithRender(ctx, render)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
