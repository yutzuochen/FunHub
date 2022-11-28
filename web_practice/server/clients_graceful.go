package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var Concurrency = 5
var PlayTime = 10
var Bet = "1"

type Client struct {
	id             int
	totalRevenue   int
	totalPlayTimes int
	rtp            int
	bet            int
}

func main() {
	clientArr := []Client{}
	ch := make(chan Client)
	for i := 1; i <= Concurrency; i++ {
		routineID := i
		go func() {
			cli := Client{id: i}

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
					Timeout:   30 * time.Second,
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

				bet, err := strconv.Atoi(Bet)
				if err != nil {
					fmt.Println(err)
				}
				cli.bet += bet
				revenue, err := strconv.Atoi(string(bds))
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("routineID", routineID, " single revenue: ", revenue)
				cli.totalRevenue += revenue
			}
			// time.Sleep(500 * time.Millisecond)

			ch <- cli
		}()

		// time.Sleep(700 * time.Millisecond) // to make seed of time different
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := client.Shutdown(ctx); err != nil {
	// 	log.Fatal("Server Shutdown: ", err)
	// }
	for n := 1; n <= Concurrency; n++ {
		clientArr = append(clientArr, <-ch)
	}
	revenue := 0
	playTimes := 0
	bet := 0
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	for _, c := range clientArr {
		fmt.Println(" client", c.id, " play", c.totalPlayTimes, " times, his revenue: ", c.totalRevenue)
		revenue += c.totalRevenue
		playTimes += c.totalPlayTimes
		bet += c.bet
	}
	fmt.Println("sum of revenue: ", revenue)
	fmt.Println("sum of bet: ", bet)
	fmt.Println("RTP: ", float64(revenue)/float64(bet))

}
