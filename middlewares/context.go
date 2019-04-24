package middlewares

import (
	"net/http"

	"github.com/MixinNetwork/mixin-wallet/session"
	"github.com/unrolled/render"
)

func Context(handler http.Handler, render *render.Render) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := session.WithRequest(r.Context(), r)
		ctx = session.WithRender(ctx, render)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
