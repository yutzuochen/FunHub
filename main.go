package main

import (
	"context"

	"net/http"
	"os"
	"os/signal"

	gc "Funhub/gamecore"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	//log輸出為json格式
	logrus.SetFormatter(&logrus.JSONFormatter{})
	//輸出設定為標準輸出(預設為stderr)
	logrus.SetOutput(os.Stdout)
	//設定要輸出的log等級
	logrus.SetLevel(logrus.DebugLevel)
}
func main() {
	logrus.Info("example log")
	logrus.WithFields(logrus.Fields{
		"trace": "trace-0001",
	}).Error("Error Log")
	// sugarLogger := zlog.InitLogger()
	// defer sugarLogger.Sync()
	ctx := context.Background()
	mux := http.NewServeMux()

	mux.HandleFunc("/play/", gc.Play)

	server := &http.Server{
		Addr:         ":1210",
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}

	// sugarLogger.Info("starting httpserver at http:localhost:1210")
	// sugarLogger.Debug("starting httpserver at http:localhost:1210_debug")
	logrus.Info("starting httpserver at http:localhost:1210")

	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	//Notify：將系統訊號轉發至channel
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//阻塞channel
	s := <-quit
	logrus.Info("Got signal: ", s)
	// If process haven't finished within 10second, we will force to shutdown the process.'
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	logrus.Info("Graceful Shutdown start ")

	if err := server.Shutdown(c); err != nil {
		logrus.Error("server.Shutdown:", err)
	}
	logrus.Info("Graceful Shutdown end ")
}