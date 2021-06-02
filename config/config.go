package config

import (
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type Config struct {
	LogPath  string         `json:"lp"`
	DBConfig DatabaseConfig `json:"db"`
	Port     int            `json:"port"` //Port at which server will run

	//We can add more config related variables here
}

type DatabaseConfig struct {
	DatabaseType     string `json:"dt"`
	ConnectionString string `json:"cs"`
}

var (
	ConfigObj Config
	Logger    *zap.Logger
)

const (
	cONFIG_PATH = "config/config.json"
)

func InitConfig() (err error) {
	//Now loading config params for config file
	cObj := new(Config)
	f, err := os.Open(cONFIG_PATH)
	if err != nil {
		fmt.Println("Error in open config file at path ", cONFIG_PATH)
		return err
	}
	jEnc := json.NewDecoder(f)
	err = jEnc.Decode(cObj)
	if err != nil {
		fmt.Println("Error in encoding config string "+
			"to config obj; error is ", err.Error())
		return err
	} else {
		ConfigObj = *cObj
	}

	//Setting up zap logger to log in file
	if err := setUpLogger(); err != nil {
		fmt.Println("Error in configuring logger; error is ", err.Error())
		return err
	}

	return nil
}

func setUpLogger() (err error) {
	lConfig := zap.NewProductionConfig()
	//Enabling file based logging

	lConfig.OutputPaths = []string{"stdout", ConfigObj.LogPath}

	if l, err := lConfig.Build(); err != nil {
		fmt.Println("Error in building zap logger ; error is ", err.Error())
		return err
	} else {
		Logger = l
		return nil
	}
}
