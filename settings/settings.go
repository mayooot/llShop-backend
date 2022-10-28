package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 全局变量，用来保存程序的所有配置信息
var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Port         int    `mapstructure:"port"`
	Version      string `mapstructure:"version"`
	StartTime    string `mapstructure:"start_time"`
	MachineId    int64  `mapstructure:"machine_id"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
	*UserConfig  `mapstructure:"user"`
	*Aliyun      `mapstructure:"aliyun"`
	*RabbitMQ    `mapstructure:"rabbitmq"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Dbname       string `mapstructure:"dbname"`
	MaxOPenConns int    `mapstructure:"max_conns"`
	MaxIdelConns int    `mapstructure:"max_idel_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type UserConfig struct {
	Width      int `mapstructure:"width"`
	MinPassLen int `mapstructure:"min_pass_len"`
	MaxPassLen int `mapstructure:"max_pass_len"`
}

type RabbitMQ struct {
	Port     int    `mapstructure:"port"`
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type Aliyun struct {
	AccessKeyId      string `mapstructure:"access_key_id"`
	AccessKeySecret  string `mapstructure:"access_key_secret"`
	AccessKeyId2     string `mapstructure:"access_key_id2"`
	AccessKeySecret2 string `mapstructure:"access_key_secret2"`
	*OSSConfig       `mapstructure:"oss"`
}

type OSSConfig struct {
	Endpoint         string `mapstructure:"endpoint"`
	BucketName       string `mapstructure:"bucket_name"`
	UserAvatarPrefix string `mapstructure:"user_avatar_prefix"`
}

func Init() (err error) {
	viper.SetConfigFile("config.yaml") // 指定配置文件
	err = viper.ReadInConfig()         // 读取配置信息

	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
		return
	}
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	}
	viper.WatchConfig() // 实时监控配置文件
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 当配置文件修改后，将配置重新序列化到Conf中
		fmt.Println("配置文件已修改...")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})
	return
}
