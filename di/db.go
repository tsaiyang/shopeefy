package di

import (
	"shopeefy/internal/repository/dao"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const dbConfigKey = "mysql"

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}

	var cfg Config
	if err := viper.UnmarshalKey(dbConfigKey, &cfg); err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN))
	if err != nil {
		panic(err)
	}

	if err = dao.InitTables(db); err != nil {
		panic(err)
	}

	return db
}
