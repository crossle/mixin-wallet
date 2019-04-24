package middlewares

import (
	"net/http"
)

func Authenticate(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// header := r.Header.Get("Authorization")
		// if !strings.HasPrefix(header, "Bearer ") {
		// 	views.RenderErrorResponse(w, r, session.AuthorizationError(r.Context()))
		// 	return
		// }
		// if header[7:] != "hello" {
		// 	views.RenderErrorResponse(w, r, session.AuthorizationError(r.Context()))
		// 	return
		// }
		handler.ServeHTTP(w, r)
	})
}
