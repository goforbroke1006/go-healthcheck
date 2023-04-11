package healthcheck

import (
	"context"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/mux"
)

const DefaultAddr = "0.0.0.0:8080"

func Panel() *HttpPanel {
	if globalHealthCheck == nil {
		globalHealthCheck = &HttpPanel{}
	}
	return globalHealthCheck
}

var globalHealthCheck *HttpPanel = nil

type HttpPanel struct {
	healthy     uint32
	readinessFn []func() bool
}

func (hc *HttpPanel) SetHealthy() {
	atomic.StoreUint32(&hc.healthy, 1)
}

func (hc *HttpPanel) SetReady() {
	hc.readinessFn = append(hc.readinessFn, func() bool {
		return true
	})
}

func (hc *HttpPanel) isReady() bool {
	if len(hc.readinessFn) == 0 {
		return false
	}

	for _, fn := range hc.readinessFn {
		if !fn() {
			return false
		}
	}
	return true
}

func (hc *HttpPanel) Start(ctx context.Context, handleAddr string) {
	handler := &mux.Router{}
	srv := &http.Server{Addr: handleAddr, Handler: handler}

	go func() {
		handler.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
			if atomic.LoadUint32(&hc.healthy) > 0 {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("OK"))
			} else {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte("Fail"))
			}
		})
		handler.HandleFunc("/readyz", func(w http.ResponseWriter, req *http.Request) {
			if hc.isReady() {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("OK"))
			} else {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte("Fail"))
			}
		})

		go func() {
			_ = srv.ListenAndServe()
		}()

		<-ctx.Done()

		_ = srv.Shutdown(context.Background())
	}()
}
