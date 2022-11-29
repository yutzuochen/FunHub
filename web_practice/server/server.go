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
)

var M = map[int]int{
	1: 0,
	2: 5,
	3: 3,
	4: 8,
	5: 0,
	6: 10,
}

func main() {
	ctx := context.Background()
	mux := http.NewServeMux()
	mux.HandleFunc("/play/", play)

	server := &http.Server{
		Addr:         ":1210",
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}

	log.Println("starting httpserver at http:localhost:1210")

	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	//Notify：將系統訊號轉發至channel
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//阻塞channel
	s := <-quit
	log.Println("Got signal: ", s)
	// If process haven't finished within 10second, we will force to shutdown the process.'
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	log.Println("Graceful Shutdown start ")

	if err := server.Shutdown(c); err != nil {
		log.Println("server.Shutdown:", err)
	}
	log.Println("Graceful Shutdown end ")
}
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func play(w http.ResponseWriter, r *http.Request) {
	bet_str := strings.TrimPrefix(r.URL.Path, "/play/")
	rand.Seed(time.Now().UnixNano())
	ran := randomInt(1, 7)
	log.Println("bet, multiple, random number: ", bet_str, M[ran], ran)
	bet, err := strconv.Atoi(bet_str)
	if err != nil {
		log.Println(err)
	}

	res := strconv.Itoa(bet * M[ran])

	fmt.Fprintf(w, res)
}
