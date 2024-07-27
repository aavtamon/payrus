// PizTec Corporation, 2024. All Rights Reserved.

package main

import (
	"fmt"

	"github.com/spf13/viper"
)

var Config *viper.Viper

func InitConfig(fileName string) error {
	v := viper.New()

	initDefaults(v)

	v.SetConfigFile(fileName)
	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("config cannot be loaded - %w", err)
	}

	Config = v

	return nil
}

func initDefaults(v *viper.Viper) {
	v.SetDefault("server.web_root", "/")
	v.SetDefault("server.web_port", 8888)
}
