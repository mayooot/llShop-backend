package gen

import (
	"github.com/bwmarrin/snowflake"
	"time"
)

var node *snowflake.Node

// Init
//
//	@Description: 雪花算法初始化
//	@param startTime 起始时间，从这个时间开始算起，可以使用雪花ID可以使用69年
//	@param machineId 机器ID
//	@return err
func Init(startTime string, machineId int64) (err error) {
	var st time.Time
	// time.Parse()解析一个格式化的时间字符串并返回它代表的时间。layout定义了参考时间
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(machineId)
	return
}

// GenSnowflakeId 生成雪花算法ID
func GenSnowflakeId() int64 {
	return node.Generate().Int64()
}
