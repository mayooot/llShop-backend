package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-backend/logic"
	"shop-backend/models/dto"
	"strconv"
)

// UserReceiverAddressListHandler 获取地址表中的所有地址，用于用户选择
// @Summary 获取地址表中的所有地址，用于用户选择
// @Description 前端需要携带Token用来鉴权，鉴权通过后返回数据库中所有地址信息
// @Tags 收货地址相关接口
// @Produce json
// @param Authorization header string true "Bearer AToken&RToken"
// @Router /user/receiveraddress/list [get]
func UserReceiverAddressListHandler(c *gin.Context) {
	pcdList, err := logic.GetAllAddress()
	if err != nil {
		zap.L().Error("获取所有收货地址失败", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, pcdList)
}

// UserReceiverAddressAddHandler 添加一条用户的收货地址
// @Summary 添加用户收货地址
// @Description 前端需传递Token信息，并封装成用户收货地址结构体。不需要传递主键ID
// @Tags 收货地址相关接口
// @Produce json
// @param Authorization header string true "Bearer AToken&RToken"
// @Param cartProductList body dto.ReceiverAddress true "用户收货地址信息结构体"
// @Router /user/receiveraddress/add [post]
func UserReceiverAddressAddHandler(c *gin.Context) {
	address := new(dto.ReceiverAddress)
	if err := c.ShouldBindJSON(address); err != nil {
		zap.L().Error("新增用户收货地址接口，传递参数有误", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	err := logic.AddReceiverAddress(address, c.GetInt64("uid"))
	if err != nil {
		zap.L().Error("新增用户收货地址接口，添加失败", zap.Error(err))
		ResponseError(c, CodeAddReceiverAddressFailed)
		return
	}
	ResponseSuccessWithMsg(c, "添加成功🎴", nil)
}

// UserReceiverAddressUpdateHandler 修改用户的一条收货地址信息
// @Summary 修改用户的一条收货地址信息
// @Description 前端需传递Token信息，并封装成用户收货地址结构体。需要传递主键ID
// @Tags 收货地址相关接口
// @Produce json
// @param Authorization header string true "Bearer AToken&RToken"
// @Param cartProductList body dto.ReceiverAddress true "用户收货地址信息结构体"
// @Router /user/receiveraddress/update [put]
func UserReceiverAddressUpdateHandler(c *gin.Context) {
	address := new(dto.ReceiverAddress)
	if err := c.ShouldBindJSON(address); err != nil {
		zap.L().Error("修改用户收货地址接口，传递参数有误", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	err := logic.UpdateReceiverAddress(address, c.GetInt64("uid"))
	if err != nil {
		zap.L().Error("修改用户收货地址接口，修改失败", zap.Error(err))
		ResponseError(c, CodeUpdateReceiverAddressFailed)
		return
	}
	ResponseSuccessWithMsg(c, "修改成功🐗", nil)
}

// UserReceiverAddressDeleteHandler 删除用户的一条收货地址信息
// @Summary 删除用户的一条收货地址信息
// @Description 前端需传递Token信息，只需要传递主键ID即可
// @Tags 收货地址相关接口
// @Produce json
// @param Authorization header string true "Bearer AToken&RToken"
// @Param id path string true "收货地址信息主键ID"
// @Router /user/receiveraddress/delete/{id} [delete]
func UserReceiverAddressDeleteHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if idStr == "" || err != nil {
		zap.L().Error("删除用户收货地址接口，传递参数有误")
		ResponseError(c, CodeInvalidParams)
		return
	}
	if err = logic.DelReceiverAddress(id, c.GetInt64("uid")); err != nil {
		zap.L().Error("删除用户收货地址接口，删除失败", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccessWithMsg(c, "删除成功🦌", nil)
}

// UserReceiverAddressPersonHandler 获取用户所有的收货地址
// @Summary 获取用户所有的收货地址
// @Description 前端需要携带Token用来鉴权，鉴权通过后返回用户所有的收货地址
// @Tags 收货地址相关接口
// @Produce json
// @param Authorization header string true "Bearer AToken&RToken"
// @Router /user/receiveraddress/my [get]
func UserReceiverAddressPersonHandler(c *gin.Context) {
	data, err := logic.GetPersonAllAddress(c.GetInt64("uid"))
	if err != nil {
		zap.L().Error("获取用户所有的收货地址失败", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, data)
}
