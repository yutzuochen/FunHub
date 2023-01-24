package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var M = map[int]int{
	1: 0,
	2: 5,
	3: 3,
	4: 8,
	5: 0,
	6: 10,
}

var sugarLogger *zap.SugaredLogger

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.MessageKey = "message"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func InitLogger() {
	fmt.Println("InitLogger")
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
	sugarLogger.Debug("Initlogger")
}
func main() {
	InitLogger()
	defer sugarLogger.Sync()
	ctx := context.Background()
	mux := http.NewServeMux()
	mux.HandleFunc("/play/", play)

	server := &http.Server{
		Addr:         ":1210",
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}

	sugarLogger.Info("starting httpserver at http:localhost:1210")
	sugarLogger.Debug("starting httpserver at http:localhost:1210_debug")

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
	sugarLogger.Info("Got signal: ", s)
	// If process haven't finished within 10second, we will force to shutdown the process.'
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	sugarLogger.Info("Graceful Shutdown start ")

	if err := server.Shutdown(c); err != nil {
		sugarLogger.Error("server.Shutdown:", err)
	}
	sugarLogger.Info("Graceful Shutdown end ")
}
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func play(w http.ResponseWriter, r *http.Request) {
	bet_str := strings.TrimPrefix(r.URL.Path, "/play/")
	rand.Seed(time.Now().UnixNano())
	ran := randomInt(1, 7)
	sugarLogger.Debug("bet, multiple, random number: ", bet_str, M[ran], ran)
	bet, err := strconv.Atoi(bet_str)
	if err != nil {
		log.Fatal(err)
	}

	res := strconv.Itoa(bet * M[ran])

	fmt.Fprintf(w, res)
}
