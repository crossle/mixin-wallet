package session

import (
	"context"
	"net/http"

	"github.com/crossle/mixin-wallet/durable"
	"github.com/golang-jwt/jwt"
	"github.com/unrolled/render"
)

type contextValueKey int

const (
	keyRequest           contextValueKey = 0
	keyDatabase          contextValueKey = 1
	keyLogger            contextValueKey = 2
	keyRender            contextValueKey = 3
	keyLimiter           contextValueKey = 4
	keyRemoteAddress     contextValueKey = 11
	keyAuthorizationInfo contextValueKey = 12
	keyRequestBody       contextValueKey = 13
)

func Database(ctx context.Context) *durable.Database {
	v, _ := ctx.Value(keyDatabase).(*durable.Database)
	return v
}

func Render(ctx context.Context) *render.Render {
	v, _ := ctx.Value(keyRender).(*render.Render)
	return v
}

func Request(ctx context.Context) *http.Request {
	v, _ := ctx.Value(keyRequest).(*http.Request)
	return v
}

func RemoteAddress(ctx context.Context) string {
	v, _ := ctx.Value(keyRemoteAddress).(string)
	return v
}

func RequestBody(ctx context.Context) string {
	v, _ := ctx.Value(keyRequestBody).(string)
	return v
}

func AuthorizationInfo(ctx context.Context) jwt.MapClaims {
	v, _ := ctx.Value(keyAuthorizationInfo).(jwt.MapClaims)
	return v
}

func WithDatabase(ctx context.Context, database *durable.Database) context.Context {
	return context.WithValue(ctx, keyDatabase, database)
}

func WithRender(ctx context.Context, render *render.Render) context.Context {
	return context.WithValue(ctx, keyRender, render)
}

func WithRequest(ctx context.Context, r *http.Request) context.Context {
	rCopy := new(http.Request)
	*rCopy = *r
	return context.WithValue(ctx, keyRequest, rCopy)
}

func WithRemoteAddress(ctx context.Context, remoteAddr string) context.Context {
	return context.WithValue(ctx, keyRemoteAddress, remoteAddr)
}

func WithRequestBody(ctx context.Context, body string) context.Context {
	return context.WithValue(ctx, keyRequestBody, body)
}

func WithAuthorizationInfo(ctx context.Context, tokenString string) context.Context {
	var value jwt.MapClaims = make(jwt.MapClaims)
	jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			value = claims
		}
		return nil, nil
	})
	return context.WithValue(ctx, keyAuthorizationInfo, value)
}
