package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

const EnvPrefix = "CHATGPT"
const AppName = "chatgpt-cli"
const CommandName = "chatgpt"
const Version = "0.0.0"
const ConfigFileName = "config"
const ConfigExt = "yaml"

const (
	ConfigApiKey     = "apikey"
	ConfigOrganizeID = "org"
)

func setDefaultConfig() {
	viper.SetDefault(ConfigApiKey, "")
	viper.SetDefault(ConfigOrganizeID, "")
}

func emptyCompleter(prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}
func initConfig(configFilePathRest ...string) error {
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	configPath := filepath.Join(configFolder, AppName)
	if len(configFilePathRest) > 0 && len(configFilePathRest[0]) > 0 {
		viper.SetConfigFile(configFilePathRest[0])
	} else {
		viper.AddConfigPath(configPath)
		viper.SetConfigName(ConfigFileName)
		viper.SetConfigType(ConfigExt)
	}

	setDefaultConfig()
	// return viper.ReadInConfig()

	err := viper.ReadInConfig()
	apiKey := viper.GetString(ConfigApiKey)
	orgID := viper.GetString(ConfigOrganizeID)
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		configFile := filepath.Join(configPath, ConfigFileName+"."+ConfigExt)
		fmt.Printf("Cannot find configuration file(%s)\n", configFile)
		for apiKey == "" {
			apiKey = strings.TrimSpace(prompt.Input("? Please input your openai api key: ", emptyCompleter))
		}
		viper.Set(ConfigApiKey, apiKey)
		for orgID == "" {
			orgID = strings.TrimSpace(prompt.Input("? Please input your organize id:  ", emptyCompleter))
		}
		viper.Set(ConfigOrganizeID, orgID)
		err = os.MkdirAll(configPath, os.FileMode(0700))
		if err != nil {
			return err
		}
		return viper.WriteConfigAs(configFile)
	}
	return nil
}
