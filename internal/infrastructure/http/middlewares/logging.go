package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

var logger = slog.New(tint.NewHandler(os.Stdout, &tint.Options{
	TimeFormat: "02-Jan-2006 15:04 PM -0700",
}))

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusResponseWriter{ResponseWriter: w}
		t0 := time.Now()
		next.ServeHTTP(sw, r)
		if sw.status >= 500 {
			logger.Error(fmt.Sprintf("\"%s %s %s\"", r.Method, r.RequestURI, r.Proto), "remote_addr", r.RemoteAddr, "code", sw.status, "time", time.Since(t0))
		} else if sw.status >= 400 {
			logger.Warn(fmt.Sprintf("\"%s %s %s\"", r.Method, r.RequestURI, r.Proto), "remote_addr", r.RemoteAddr, "code", sw.status, "time", time.Since(t0))
		} else {
			logger.Info(fmt.Sprintf("\"%s %s %s\"", r.Method, r.RequestURI, r.Proto), "remote_addr", r.RemoteAddr, "code", sw.status, "time", time.Since(t0))
		}
	})
}
