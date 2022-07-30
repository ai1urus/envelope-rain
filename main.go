package main

import (
	"envelope-rain/router"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()
	router.InitService()
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	r.POST("/snatch", router.SnatchHandler)
	r.POST("/open", router.OpenHandler)
	r.POST("/get_wallet_list", router.WalletListHandler)

	r.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	router.StopService()
}
