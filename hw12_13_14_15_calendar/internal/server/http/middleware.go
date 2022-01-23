package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler, logg Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, r)
		latency := fmt.Sprintf("%dms", time.Since(t).Milliseconds())
		ret := fmt.Sprintf("middleware: %s %s %s %s %d %s %s\n",
			r.RemoteAddr, r.Method, r.RequestURI, r.Proto, 200, latency, r.UserAgent())
		fmt.Fprintf(w, "<p>%s</p>", ret)
		logg.Info(ret)
	})
}
