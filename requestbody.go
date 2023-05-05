package requestbody

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "github.com/caddyserver/caddy/v2"
)

type Middleware struct {
    Next    http.Handler
}

func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    buffer := &bytes.Buffer{}
    tee := io.TeeReader(r.Body, buffer)

    var bodyJSON interface{}
    if err := json.NewDecoder(tee).Decode(&bodyJSON); err != nil && err != io.EOF {
        caddy.Log().Error("Error decoding request body", zap.Error(err))
    }

    // Restore the original request body before passing to the next handler.
    r.Body = ioutil.NopCloser(bytes.NewReader(buffer.Bytes()))

    logger := getLogger(r.Context())
    logger.Info("Request Body", zap.Any("request", bodyJSON))

    m.Next.ServeHTTP(w, r)
}

func getLogger(ctx context.Context) *zap.Logger {
    return ctx.Value(caddy.LogContextKey).(*zap.Logger)
}

func New() *Middleware {
    return &Middleware{}
}

func (mw *Middleware) CaddyModule() caddy.ModuleInfo {
    return caddy.ModuleInfo{
        ID:  "http.handlers.requestbody",
        New: func() caddy.Module { return new(Middleware) },
    }
}

func (mw *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
    mw.Next = next
    mw.ServeHTTP(w, r)
    return nil
}
