package controller

import (
	"github.com/gin-gonic/gin"
	"log"
)

func SecKillAllSkuHandler(c *gin.Context) {
	log.Print("秒杀商品成功")
	ResponseSuccess(c, nil)
}
