package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	router := gin.Default()
	router.GET("/DoGetByQueryString", func(c *gin.Context) {
		time.Sleep(20 * time.Second)
		p1 := c.DefaultQuery("param1", "Default")
		p2 := c.Query("param2")
		c.JSON(http.StatusOK, gin.H{"param1": p1, "param2": p2})
	})

	//原本是用router.Run()，要使用net/http套件的shutdown的話，需要使用原生的ListenAndServe
	srv := &http.Server{
		Addr:    ":8787",
		Handler: router,
	}
	//新增一個channel，type是os.Signal
	ch := make(chan os.Signal, 1)
	//call goroutine啟動http server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("SERVER GG惹:", err)
		}
	}()
	//Notify：將系統訊號轉發至channel
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	//阻塞channel
	<-ch

	//收到關機訊號時做底下的流程
	fmt.Println("Graceful Shutdown start - 1")
	//透過context.WithTimeout產生一個新的子context，它的特性是有生命週期，這邊是設定10秒
	//只要超過10秒就會自動發出Done()的訊息
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	fmt.Println("Graceful Shutdown start - 2")
	//使用net/http的shutdown進行關閉http server，參數是上面產生的子context，會有生命週期10秒，
	//所以10秒內要把request全都消化掉，如果超時一樣會強制關閉，所以如果http server要處理的是
	//需要花n秒才能處理的request就要把timeout時間拉長一點
	if err := srv.Shutdown(c); err != nil {
		log.Println("srv.Shutdown:", err)
	}
	//使用select去阻塞主線程，當子context發出Done()的訊號才繼續向下走
	select {
	case <-c.Done():
		fmt.Println("Graceful Shutdown start - 3")
		close(ch)
	}
	fmt.Println("Graceful Shutdown end ")
}
