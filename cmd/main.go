package main

import (
	"fmt"
	"go.uber.org/zap"
	"minesweeper_api"
	"minesweeper_api/api"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)


const svcName = "minesweep"

func main() {

	logger, _ := zap.NewProduction()
	svc := minesweeper_api.NewService(logger)

	errs := make(chan error, 2)

	go startHTTPServer(svc, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err := <-errs
	logger.Error(fmt.Sprintf("Post service terminated: %s", err))
}

func startHTTPServer(svc minesweeper_api.MineSweepService, logger *zap.Logger, errs chan error) {
	port := fmt.Sprintf(":%d", 8080)
	logger.Info(fmt.Sprintf("Post service started using http on port %d", 8080))
	errs <- http.ListenAndServe(port, api.MakeHandler(svcName, svc))
}
