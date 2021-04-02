package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Service struct {
	ctx           context.Context
	hasContext    bool
	listenAddress string
	checkFunc     func() int
}

type Option func(svc *Service)

func (s *Service) Start() {
	router := http.NewServeMux()
	router.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(s.checkFunc())
	}))

	srv := &http.Server{
		Addr:         s.listenAddress,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("[HEALTHCHECK] web service on %s exited: %v\n", s.listenAddress, err)
		}
	}()

	// Wait for the context to finish and shut down the health check service
	if s.hasContext {
		go func() {
			<-s.ctx.Done()

			srv.SetKeepAlivesEnabled(false)
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := srv.Shutdown(shutdownCtx); err != nil {
				fmt.Println("[HEALTHCHECK] could not gracefully shutdown the service:", err)
			}

			fmt.Println("[HEALTHCHECK] web service stopped")
		}()
	}
}

func New(opts ...Option) *Service {
	svc := &Service{
		ctx:           context.Background(),
		listenAddress: ":11223",
		checkFunc: func() int {
			return http.StatusOK
		},
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

func ListenOn(addr string) Option {
	return func(svc *Service) {
		svc.listenAddress = addr
	}
}

func WithContext(ctx context.Context) Option {
	return func(svc *Service) {
		svc.hasContext = true
		svc.ctx = ctx
	}
}

func WithCustomCheck(fnc func() int) Option {
	return func(svc *Service) {
		if fnc != nil {
			svc.checkFunc = fnc
		}
	}
}
