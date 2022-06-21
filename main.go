package main

import (
	"envelope-rain/router"
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	router.InitService()
	// go func() {
	// 	log.Println(http.ListenAndServe(":6060", nil))
	// }()

	r.POST("/snatch", router.SnatchHandler)
	r.POST("/open", router.OpenHandler)
	r.POST("/get_wallet_list", router.WalletListHandler)

	r.Run()
}
