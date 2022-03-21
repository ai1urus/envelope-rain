// package main

// import (
// 	"envelope-rain/database"

// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// func main() {
// 	dsn := "root:852196@tcp(175.27.248.24:3306)/envelope?charset=utf8mb4&parseTime=True&loc=Local"
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		panic("failed to connect database")
// 	}

// 	db.AutoMigrate(&database.User{})

// 	user := &database.User{
// 		Uid:       2,
// 		Amount:    8000,
// 		Cur_count: 5,
// 	}

// 	db.Create(&user)
// }

package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()
	// router
	r.POST("/snatch", SnatchHandler)
	r.POST("/open", OpenHandler)
	r.POST("/get_wallet_list", WalletListHandler)

	r.Run()
}

func SnatchHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.WithFields(log.Fields{
		"uid": uid,
	}).Info("User snatched")

	// logic start
	envelope_id := 123
	max_count := 5
	cur_count := 3
	// logic end

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"envelope_id": envelope_id,
			"max_count":   max_count,
			"cur_count":   cur_count,
		},
	})
}

func OpenHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	envelope_id, _ := c.GetPostForm("envelope_id")

	log.Info("envelope %d opened by %d", envelope_id, uid)

	// logic start
	value := 50
	// logic end

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"value": value,
		},
	})
}

func WalletListHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	log.Info("query %d's wallet", uid)

	// logic start
	envelopes := []gin.H{
		{
			"envelope_id": 123,
			"value":       50,
			"opened":      true,
			"snatch_time": 1634551711,
		},
		{
			"envelope_id": 234,
			"opened":      false,
			"snatch_time": 1634551711,
		},
	}
	amount := 50
	// logic end

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"amount":        amount,
			"envelope_list": envelopes,
		},
	})
}
