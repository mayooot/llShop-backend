package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-backend/logic"
	"shop-backend/models/dto"
	"shop-backend/rabbitmq"
	"strconv"
)

// SecKillAllSkuHandler 获取所有正在秒杀的商品
// @Summary 获取所有正在秒杀的商品
// @Description 前端不需要携带Token
// @Tags 秒杀相关接口
// @Produce json
// @Router /seckill/sku/list [get]
func SecKillAllSkuHandler(c *gin.Context) {
	data, err := logic.GetAllSecKillSku()
	if err != nil {
		ResponseBadError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, data)
}

// SecKillBuyHandler 秒杀商品接口
func SecKillBuyHandler(c *gin.Context) {
	product := new(dto.SecKillProduct)
	if err := c.ShouldBindJSON(product); err != nil {
		zap.L().Error("秒杀商品接口，传递参数错误")
		ResponseBadError(c, CodeServeBusy)
		return
	}

	skuID, err := strconv.ParseInt(product.SkuID, 10, 64)
	if err != nil {
		zap.L().Error("秒杀商品接口，商品skuID转为int64错误")
		ResponseError(c, CodeServeBusy)
		return
	}

	// 发布到MQ中
	err = rabbitmq.SendSecKillReqMess2MQ(&dto.SecKillMQ{
		SkuID: skuID,
		UID:   c.GetInt64("uid"),
	})
	if err != nil {
		// 发布到MQ失败
		ResponseBadError(c, CodeSecKillFinished)
		return
	}

	ResponseSuccessWithMsg(c, "正在秒杀中，请在我的订单查看是否秒杀成功🍔", nil)
}
