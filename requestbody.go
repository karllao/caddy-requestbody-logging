package requestbody

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
    "go.uber.org/zap/zapcore"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(Middleware{})
}

type Middleware struct{}

func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.requestbody",
		New: func() caddy.Module { return new(Middleware) },
	}
}

func getLogger(ctx context.Context) *zap.Logger {
    return ctx.Value(caddy.LogContextKey).(*zap.Logger)
}

func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

    logger := getLogger(r.Context())
    logger.Info("Request Body", zap.Any("request", string(bodyBytes)))
	
	return next.ServeHTTP(w, r)
}