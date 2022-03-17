package quark

import (
	"log"
	"net"
	"net/http"
	"strings"
)

type Middleware func(http.Handler) http.Handler

func Chain(m ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLogResponseWriter(w http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{w, http.StatusOK}
}

func (lw *logResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}

func AccessLog(next http.Handler) http.Handler {
	// 第三题
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logW := NewLogResponseWriter(w)
		next.ServeHTTP(logW, r)
		clientIP := r.Header.Get("X-Forwarded-For")
		clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
		if clientIP == "" {
			clientIP = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
		}
		if clientIP != "" {
			log.Printf("remote: %s | code: %d | %s | %s", clientIP, logW.statusCode, r.Method, r.URL)
			return
		}
		if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
			log.Printf("remote: %s | %d | %s | %s", ip, logW.statusCode, r.Method, r.URL)

		}
	})
}
