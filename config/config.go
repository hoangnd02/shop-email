package config

import (
	"github.com/gookit/goutil/fsutil"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB
var Mailer struct {
	Events []*struct {
		Key      string `yaml:"key"`
		Name     string `yaml:"name"`
		Template string `yaml:"template"`
	}
}

func InnitConfig() {
	dsn := "host=localhost user=root password=123456 dbname=auth port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	Database = db

	template := fsutil.MustReadFile("config/mailer.yml")
	yaml.Unmarshal(template, &Mailer)
}
