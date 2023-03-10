package ruleconfig

import (
	"fmt"
	"github.com/spf13/viper"
)

var GlobalConfig Config

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	DB   int    `mapstructure:"db"`
}

type Config struct {
	Port         string      `mapstructure:"port"`
	DefaultDb    DBConfig    `mapstructure:"default_db"`
	DefaultRedis RedisConfig `mapstructure:"default_redis"`
}

func (dbConfig DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
	)
}

func (redisCfg RedisConfig) DSN() string {
	return fmt.Sprintf("%s:%d",
		redisCfg.Host,
		redisCfg.Port,
	)
}

func init() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/ruleconfig")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&GlobalConfig)
	fmt.Println(GlobalConfig)
}
