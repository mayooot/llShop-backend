package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"shop-backend/settings"
)

var db *gorm.DB
var sqlDB *sql.DB

// Init 初始化gorm和mysql
func Init(cfg *settings.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Dbname)
	// 初始化gorm
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// 初始化失败
		zap.L().Error("gorm.Open(mysql.Open(dsn), &gorm.Config{}) failed", zap.Error(err))
		return
	}
	// 获取sql.DB
	sqlDB, _ = db.DB()
	if err = sqlDB.Ping(); err != nil {
		// 确认数据库是否连接正确
		zap.L().Error("sqlDB.Ping() failed", zap.Error(err))
		return
	}
	// 设置数据库连接池最大连接数
	sqlDB.SetMaxOpenConns(cfg.MaxOPenConns)
	// 设置数据库连接池空闲连接的上限
	sqlDB.SetMaxIdleConns(cfg.MaxIdelConns)
	return
}

// Close 关键数据库连接
func Close() {
	_ = sqlDB.Close()
}
