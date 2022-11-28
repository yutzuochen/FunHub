package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
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
	// 创建路由器
	mux := http.NewServeMux()
	// 设置路由规则
	mux.HandleFunc("/play/", play)

	// 创建服务器
	server := &http.Server{
		Addr:         ":1210",
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}

	// 监听端口并提供服务
	log.Println("starting httpserver at http:localhost:1210")
	log.Fatal(server.ListenAndServe())
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
