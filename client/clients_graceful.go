package main

import (
	"context"
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

var Concurrency = 3
var PlayTimes = 1000
var Bet = "1"

type Client struct {
	id             int
	totalRevenue   int
	totalPlayTimes int
	rtp            int
	bet            int
}

func withContextFunc(ctx context.Context, f func()) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(c)

		select {
		case <-ctx.Done():
			log.Println("[withContextFunc]case <-ctx.Done()")
		case <-c:
			log.Println("[withContextFunc]case <-c:")
			cancel()
			f()
		}
	}()
	return ctx
}

func virtualCliPlay(ctx context.Context, routineID int, clientCh chan Client) {
	cli := Client{id: routineID}
ForEnd:
	for {
		select {
		case <-ctx.Done():
			log.Println("close the worker", cli.id)
			// send result of the client
			clientCh <- cli
			return
		default:
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

			bet, err := strconv.Atoi(Bet)
			if err != nil {
				log.Fatal(err)
				// fmt.Println(err)
			}

			revenue, err := strconv.Atoi(string(bds))
			if err != nil {
				log.Fatal(err)
			}
			log.Println("routineID", routineID, " single revenue: ", revenue)
			cli.totalPlayTimes++
			cli.bet += bet
			cli.totalRevenue += revenue

			if cli.totalPlayTimes >= PlayTimes {
				break ForEnd
			}
		}
	}
	clientCh <- cli

}

func main() {
	clientArr := []Client{}
	clientCh := make(chan Client)
	finished := make(chan bool)

	ctx := withContextFunc(context.Background(), func() {
		log.Println("cancel from ctrl+c event")
		close(finished)
	})

	for i := 1; i <= Concurrency; i++ {
		go virtualCliPlay(ctx, i, clientCh)
		time.Sleep(50 * time.Millisecond) // to make seed of time different
	}
	// <-finished
	for n := 1; n <= Concurrency; n++ {
		clientArr = append(clientArr, <-clientCh)
	}
	revenue := 0
	playTimes := 0
	bet := 0

	for _, c := range clientArr {
		fmt.Println(" client", c.id, " play", c.totalPlayTimes, " times, his revenue: ", c.totalRevenue)
		revenue += c.totalRevenue
		playTimes += c.totalPlayTimes
		bet += c.bet
	}
	fmt.Println("sum of play times: ", playTimes)
	fmt.Println("sum of revenue: ", revenue)
	fmt.Println("sum of bet: ", bet)
	fmt.Println("RTP: ", float64(revenue)/float64(bet))

}
