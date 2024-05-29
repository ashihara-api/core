package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// UserAgent is the key
const UserAgent = "User-Agent"

// Logger is a middleware that logs the start and end of each request
func Logger(l *slog.Logger, ignoreUAs []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				if !ignoreUserAgents(r.Header.Get(UserAgent), ignoreUAs) {
					lat := time.Since(t1)
					l.Info("response",
						slog.Group("request",
							slog.Int("return-status", ww.Status()),
							slog.String("http-method", r.Method),
							slog.Group("headers",
								slog.String("Content-Type", r.Header.Get("Content-Type")),
								slog.String("Content-Length", r.Header.Get("Content-Length")),
								slog.String("User-Agent", r.Header.Get(UserAgent)),
								slog.String("Server", r.Header.Get("Server")),
								slog.String("Via", r.Header.Get("Via")),
								slog.String("Accept", r.Header.Get("Accept")),
								slog.String("Authorization", r.Header.Get("Authorization")),
							),
						),
						// Other data
						slog.String("X-FORWARDED-FOR", r.Header.Get("X-FORWARDED-FOR")),
						slog.String("Remote Addr", r.RemoteAddr),
						slog.String("Proto", r.Proto),
						slog.String("Path", r.URL.Path),
						slog.Duration("lat", lat),
						slog.Int("size", ww.BytesWritten()),
						slog.String("reqId", middleware.GetReqID(r.Context())),
					)
				}
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

// ignoreUserAgents checks whether the user agent of the request is a target for suppressing log output.
//
// For example, it is used to suppress log output during health checks.
//
//	kubernetes: kube-probe
//	AWS: ELB-HealthChecker
//	GCP: GoogleHC
//	Azure: Edge Health Probe & Load Balancer Agent
func ignoreUserAgents(ua string, ignores []string) bool {
	for _, u := range ignores {
		if strings.HasPrefix(ua, u) {
			return true
		}
	}
	return false
}
