package handle

import (
	"fmt"
	"indexof/logger"
	"net"
	"net/http"
	"strings"
	"time"
)

func New() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", IndexOfHandler)
	return middlewareLogger(mux)
}

type httpResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *httpResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func middlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path
		raw := r.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}
		var clientIP string
		clientIP, _, _ = net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
		if clientIP == "" {
			remoteIPHeaders := []string{"X-Forwarded-For", "X-Real-IP"}
			for _, headerName := range remoteIPHeaders {
				if headerIp := r.Header.Get(headerName); headerIp != "" {
					ips := strings.Split(headerIp, ",")
					clientIP = strings.TrimSpace(ips[0])
					if clientIP != "" {
						break
					}
				}
			}
		}
		responseWriter := &httpResponseWriter{ResponseWriter: w}
		// 调用实际的处理器
		next.ServeHTTP(responseWriter, r)
		endTime := time.Now()
		latency := endTime.Sub(start)
		if latency > time.Minute {
			latency = latency.Truncate(time.Second)
		}

		logger.Logger(fmt.Sprintf("[IndexOf] %s | %13v | %s | %d | %s  %s \n",
			endTime.Format("2006/01/02-15:04:05"),
			latency,
			clientIP,
			responseWriter.statusCode,
			r.Method,
			path,
		))
	})
}
