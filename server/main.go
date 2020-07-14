package main

import (
	"fmt"
	"strings"

	app "twistserver/app"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {

	// From the environment
	viper.SetEnvPrefix("TWIST_Example")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// From config file
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Warn("No configuration file was loaded")
	}
}

func main() {

	// Initializing application
	fmt.Println(viper.GetString("host.serviceHost"))
	fmt.Println(viper.GetString("host.twistHost"))
	a := app.CreateApp()

	err := a.Init()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Starting application
	err = a.Run()
	if err != nil {
		log.Fatal(err)
		return
	}

}
