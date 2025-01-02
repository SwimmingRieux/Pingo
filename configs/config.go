package configs

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	Errors                 Errors `json:"errors"`
	GoroutinesMax          int    `json:"goroutines_max"`
	PingerGoroutinesMax    int    `json:"pinger_goroutines_max"`
	PortsCountLimit        int    `json:"ports_count_limit"`
	ListenerIterationLimit int    `json:"listener_iteration_limit"`
	DomainsBigEnough       int    `json:"domains_big_enough"`
	V2                     V2     `json:"v2"`
}

type Errors struct {
	NotEnoughPortsFound    string `json:"not_enough_ports_found"`
	LoadFromLinkError      string `json:"load_from_link_error"`
	ConfigNotFound         string `json:"config_not_found"`
	GroupNotFound          string `json:"group_not_found"`
	ConfigFormatError      string `json:"config_format_error"`
	GroupCreatingError     string `json:"group_creating_error"`
	DirectoryCreatingError string `json:"directory_creating_error"`
	FileCreatingError      string `json:"file_creating_error"`
	HttpStatus             string `json:"http_status"`
	FileRemoveError        string `json:"file_remove_error"`
	WriteToFileError       string `json:"write_to_file_error"`
	ConfigRemoveError      string `json:"config_remove_error"`
	InvalidFormatter       string `json:"invalid_formatter"`
	InvalidPortSetter      string `json:"invalid_port_setter"`
}

type V2 struct {
	ConfigurationPath string `json:"config_path"`
}

func getDefaultConfig() string {
	return "./config/config.json"
}

func NewConfig() (*Configuration, error) {
	path := os.Getenv("cfgPath")
	if path == "" {
		path = getDefaultConfig()
	}
	viperConfig := viper.New()
	viperConfig.SetConfigName(path)
	viperConfig.AddConfigPath(".")
	viperConfig.AutomaticEnv()
	if err := viperConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	var DefaultConfig Configuration

	err := viperConfig.Unmarshal(&DefaultConfig)
	if err != nil {
		return nil, err
	}

	return &DefaultConfig, nil
}
