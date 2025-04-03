package main

import (
	"log"
	"shopeefy/pkg/logger"

	"github.com/spf13/viper"
)

func main() {
	logger.InitLogger()
	initConfig()

	server := InitWebServer()

	log.Fatalln(server.Run(":8080"))
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("dev")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	val := viper.Get("test.algoshop")
	log.Println(val)
}
