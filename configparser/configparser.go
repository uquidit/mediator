package configparser

import (
	"fmt"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var (
	Verbose bool
)

func ReadConf(config_name string, config any, defaults map[string]any) error {
	if Verbose {
		viper.Set("Verbose", true)
		jww.SetLogThreshold(jww.LevelTrace)
		jww.SetStdoutThreshold(jww.LevelTrace)
	}

	if config_name == "" {
		return fmt.Errorf("no configuration file name")
	} else {
		viper.SetConfigName(config_name)
	}
	viper.SetConfigType("yaml")                                  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("../.")                                  // path to look for the config file in
	viper.AddConfigPath(".")                                     // look for config in the working directory
	viper.AddConfigPath("/opt/mediator/conf")                    // path to look for the config file in
	viper.AddConfigPath("/opt/tufin/data/securechange/scripts/") // path to look for the config file in

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		return fmt.Errorf("fatal error while reading config file: %w ", err)
	}

	// Set default values
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	if err := viper.Unmarshal(config); err != nil {
		return fmt.Errorf("unable to decode into struct: %w", err)
	}
	return nil
}

func ReadConfAbsolutePath(config_name string, config any, defaults map[string]any) error {
	if config_name == "" {
		return fmt.Errorf("no configuration file name")
	} else {
		viper.SetConfigFile(config_name)
	}

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		return fmt.Errorf("fatal error while reading config file: %w ", err)
	}

	// Set default values
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	if err := viper.Unmarshal(config); err != nil {
		return fmt.Errorf("unable to decode into struct: %w", err)
	}
	return nil
}
