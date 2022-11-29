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

	// 监听端口并提供服务
	log.Println("starting httpserver at http:localhost:1210 1")

	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("starting httpserver at http:localhost:1210 2")
	quit := make(chan os.Signal, 1)
	log.Println("starting httpserver at http:localhost:1210 3")
	//Notify：將系統訊號轉發至channel
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//阻塞channel
	// <-quit
	s := <-quit
	fmt.Println("Got signal: ", s)
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	fmt.Println("Graceful Shutdown start - 1")
	//使用net/http的shutdown進行關閉http server，參數是上面產生的子context，會有生命週期10秒，
	//所以10秒內要把request全都消化掉，如果超時一樣會強制關閉，所以如果http server要處理的是
	//需要花n秒才能處理的request就要把timeout時間拉長一點
	if err := server.Shutdown(c); err != nil {
		log.Println("server.Shutdown:", err)
	}
	fmt.Println("Graceful Shutdown end ")
}
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func play(w http.ResponseWriter, r *http.Request) {
	bet_str := strings.TrimPrefix(r.URL.Path, "/play/")
	rand.Seed(time.Now().UnixNano())
	ran := randomInt(1, 7)
	fmt.Println("bet, multiple, random number: ", bet_str, M[ran], ran)
	bet, err := strconv.Atoi(bet_str)
	if err != nil {
		fmt.Println(err)
	}

	res := strconv.Itoa(bet * M[ran])

	fmt.Fprintf(w, res)
}
