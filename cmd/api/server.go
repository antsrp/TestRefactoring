package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func setConfig(h http.Handler, host, port string) *http.Server {
	//addr := net.JoinHostPort("", "5000")
	addr := net.JoinHostPort(host, port)

	return &http.Server{Addr: addr, Handler: h}
}

func startServer(logger *zap.Logger, r http.Handler, host, port string, sigquit chan os.Signal) {
	stopAppCh := make(chan struct{})

	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)

	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGTERM)

	srv := setConfig(r, host, port)

	go func() {
		<-sigquit
		fmt.Println("server to shutdown")

		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Sugar().Fatalf("could not shutdown server: %s", err)
		}
		stopAppCh <- struct{}{}
	}()

	if err := srv.ListenAndServe(); err != nil {
		logger.Sugar().Infof("can't listen and serve server: %s", err)
		//sigquit <- syscall.SIGTERM
	}
}

/*
func handleCloser(logger *zap.Logger, resource string, closer io.Closer) {
	if err := closer.Close(); err != nil {
		logger.Sugar().Errorf("Can't close %q: %s", resource, err)
	}
}*/
