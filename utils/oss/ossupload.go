package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.uber.org/zap"
	"io"
	"shop-backend/settings"
	"shop-backend/utils/concatstr"
	"shop-backend/utils/gen"
	"strconv"
)

var bucket *oss.Bucket
var userAvatarPrefix string
var commonPrefix = "https://llshop-project.oss-cn-zhangjiakou.aliyuncs.com/"

// Init 初始化阿里云OSS服务
func Init(cfg *settings.Aliyun) error {
	client, err := oss.New(
		cfg.Endpoint,
		cfg.AccessKeyId,
		cfg.AccessKeySecret,
	)
	userAvatarPrefix = cfg.UserAvatarPrefix
	if err != nil {
		zap.L().Error("init Aliyun OSS failed", zap.Error(err))
		return err
	}
	bucket, err = client.Bucket(cfg.BucketName)
	return err
}

// UploadPic 上传文件到阿里云服务器
func UploadPic(file io.Reader) (string, error) {
	// 雪花算法生成全局唯一图片名称
	id := gen.GenSnowflakeId()
	// 将int64转换为字符串
	idStr := strconv.FormatInt(id, 10)
	// 生成文件名
	fileName := concatstr.ConcatString(userAvatarPrefix, idStr, ".png")
	// 上传文件
	if err := bucket.PutObject(fileName, file); err != nil {
		// 上传失败
		return "", err
	}
	return concatstr.ConcatString(commonPrefix, fileName), nil
}
