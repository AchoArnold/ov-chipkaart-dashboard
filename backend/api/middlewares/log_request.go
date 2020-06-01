package middlewares

import (
	"net/http"
	"runtime/debug"
	"time"

	internalTime "github.com/NdoleStudio/ov-chipkaart-dashboard/backend/shared/time"

	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/shared/logger"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

// Status returns the status of the response writer
func (rw *responseWriter) Status() int {
	return rw.status
}

// WhiteHeader writes the header
func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

// LogRequest logs the incoming HTTP request & its duration.
func (middleware Client) LogRequest(logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_ = logger.Log(
						"err", err,
						"trace", debug.Stack(),
					)
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			_ = logger.Log(
				"at", time.Now().Format(internalTime.DefaultFormat),
				"status", wrapped.status,
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration", time.Since(start),
				"ip", r.RemoteAddr,
			)
		}

		return http.HandlerFunc(fn)
	}
}
