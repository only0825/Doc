package common

import (
	"fmt"
	"go-data/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql() (*gorm.DB, error) {
	dbinfo := configs.Conf.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		dbinfo.Username, dbinfo.Password, dbinfo.Host, dbinfo.Port, dbinfo.Db, dbinfo.Charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
