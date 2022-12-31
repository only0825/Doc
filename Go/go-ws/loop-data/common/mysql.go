package common

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"loop-data/configs"
	"loop-data/model"
)

func InitMysql() error {
	dbinfo := configs.Conf.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		dbinfo.Username, dbinfo.Password, dbinfo.Host, dbinfo.Port, dbinfo.Db, dbinfo.Charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	model.DB = db
	return nil
}
