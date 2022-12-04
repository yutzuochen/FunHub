package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

var Concurrency = 3
var Play_time = 2
var Bet = "10"

func main() {

	arr := []int{}
	for i := 1; i <= Concurrency; i++ {
		arr = append(arr, 0)
		// 创建连击池
		go func() {
			routineID := strconv.Itoa(i)
			transport := &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:          100,              // 最大空闲连接数
				IdleConnTimeout:       90 * time.Second, // 空闲超时时间
				TLSHandshakeTimeout:   10 * time.Second, // tls 握手超时时间
				ExpectContinueTimeout: 1 * time.Second,  // 100-continue状态码超时时间
			}
			// 创建客户端
			client := &http.Client{
				Transport: transport,
				Timeout:   30 * time.Second, // 没饿
			}
			url := "http://localhost:1210/play/" + Bet
			// 请求数据
			resp, err := client.Get(url)

			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			// 读取数据
			bds, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			fmt.Println("Clinent", routineID, "'s revenue: ", string(bds))
		}()
		time.Sleep(300 * time.Millisecond) // to make seed of time different
	}
	time.Sleep(1000 * time.Millisecond)
}
