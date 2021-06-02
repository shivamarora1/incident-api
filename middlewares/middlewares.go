package middlewares

import (
	"net/http"

	"example.com/incident-api/config"
	"go.uber.org/zap"
)

type handler func(w http.ResponseWriter, r *http.Request)

func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.Logger.Info("logging middleware",
			zap.String("url", r.RequestURI), zap.String("method", r.Method),
			zap.String("ip", r.RemoteAddr))
		h.ServeHTTP(w, r)
	})
}
