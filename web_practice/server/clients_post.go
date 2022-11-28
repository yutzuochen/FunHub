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
var PlayTime = 3
var Bet = "1"

type Client struct {
	id             int
	totalRevenue   int
	totalPlayTimes int
	rtp            int
}

func main() {
	clientArr := []Client{}
	// ch_revenue := make(chan int)
	// ch_playTimes := make(chan int)
	ch := make(chan Client)
	for i := 1; i <= Concurrency; i++ {
		routineID := i
		go func() {
			cli := Client{id: i}
			// totalBet := 0
			// totalRevenue := 0
			// rtp := 0
			// totalPlayTimes := 0

			for t := 1; t <= PlayTime; t++ {
				transport := &http.Transport{
					DialContext: (&net.Dialer{
						Timeout:   30 * time.Second,
						KeepAlive: 30 * time.Second,
					}).DialContext,
					MaxIdleConns:          100,
					IdleConnTimeout:       90 * time.Second,
					TLSHandshakeTimeout:   10 * time.Second,
					ExpectContinueTimeout: 1 * time.Second,
				}
				// create client side
				client := &http.Client{
					Transport: transport,
					Timeout:   30 * time.Second, // 没饿
				}
				url := "http://localhost:1210/play/" + Bet
				// request data
				resp, err := client.Get(url)

				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()

				// read data
				bds, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}
				cli.totalPlayTimes += 1

				revenue, err := strconv.Atoi(string(bds))
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("routineID", routineID, " single revenue: ", revenue)
				cli.totalRevenue += revenue
			}
			time.Sleep(500 * time.Millisecond)
			// fmt.Println("play times: ", totalPlayTimes)
			// fmt.Println("Clinent", routineID, "'s revenue: ", totalRevenue)
			// ch_revenue <- totalRevenue
			// ch_playTimes <- totalPlayTimes
			ch <- cli
		}()
		time.Sleep(700 * time.Millisecond) // to make seed of time different
	}
	// for c := 1; c <= Concurrency; c++ {
	// 	revenueArr = append(revenueArr, <-ch_revenue)
	// }
	// for c := 1; c <= Concurrency; c++ {
	// 	playTimesArr = append(playTimesArr, <-ch_playTimes)
	// }
	for n := 1; n <= Concurrency; n++ {
		clientArr = append(clientArr, <-ch)
	}
	revenue := 0
	playTimes := 0
	for c := clientArr{
		revenue += c.totalRevenue
		playTimes += c.totalPlayTimes
	}
	fmt.Println("sum of revenue: ", revenue)
	// fmt.Println("sum of playTimes: ", playTimes)
	time.Sleep(500 * time.Millisecond)
}
