package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Client *gorm.DB

func Init(options *DBOptions) {
	var err error

	Client, err = gorm.Open(mysql.Open(options.DSN))
	if err != nil {
		panic(err)
	}
}
