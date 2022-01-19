package rdb

import (
	"errors"
	"fmt"
	"github.com/ehwjh2010/viper/client/enum"
	"github.com/ehwjh2010/viper/client/setting"
	"github.com/ehwjh2010/viper/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

const defaultCreateBatchSize = 1000

var UnsupportedDBType = errors.New("invalid db type")

func InitDBWithGorm(dbConfig *setting.DB, dbType enum.DBType) (*gorm.DB, error) {

	var sqlLogger = logger.Silent
	if dbConfig.EnableRawSQL {
		sqlLogger = logger.Info
	}

	var createBatchSize = defaultCreateBatchSize
	if dbConfig.CreateBatchSize > 0 {
		createBatchSize = dbConfig.CreateBatchSize
	}

	dialector, err := getDialector(dbConfig, dbType)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		//打印SQL
		Logger: logger.Default.LogMode(sqlLogger),
		NamingStrategy: schema.NamingStrategy{
			//指定表前缀
			TablePrefix: dbConfig.TablePrefix,
			//表复数禁用
			SingularTable: dbConfig.SingularTable,
		},
		//批量操作 每批数量
		CreateBatchSize: createBatchSize,
	})

	if err != nil {
		//log.Fatalf("Connect mysql failed! err: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()

	if err != nil {
		//log.Fatalf("Access sqlDB failed! err: %v", err)
		return nil, err
	}

	// SetMaxIdleConns 设置连接池中空闲连接最大数
	sqlDB.SetMaxIdleConns(dbConfig.MaxFreeConnCount)

	// SetMaxOpenConns 设置打开数据库最大连接数
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConnCount)

	// SetConnMaxIdleTime 设置闲置连接最长存活时间
	sqlDB.SetConnMaxIdleTime(dbConfig.FreeMaxLifetime * time.Minute)

	return db, nil
}

func getDialector(dbConfig *setting.DB, dbType enum.DBType) (gorm.Dialector, error) {
	switch dbType {
	case enum.Mysql:
		dsn := fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s`,
			dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database, dbConfig.Location)
		return mysql.Open(dsn), nil

	case enum.Postgresql:
		dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s`,
			dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Database, dbConfig.Port, dbConfig.Location)

		return postgres.Open(dsn), nil
	default:
		log.Debug("only support mysql, postgresql")
		return nil, UnsupportedDBType
	}

}
