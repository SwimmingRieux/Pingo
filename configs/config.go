package configs

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	Errors                 Errors `mapstructure:"errors"`
	GoroutinesMax          int    `mapstructure:"goroutines_max"`
	PingerGoroutinesMax    int    `mapstructure:"pinger_goroutines_max"`
	PortsCountLimit        int    `mapstructure:"ports_count_limit"`
	ListenerIterationLimit int    `mapstructure:"listener_iteration_limit"`
	DomainsBigEnough       int    `mapstructure:"domains_big_enough"`
	MaxPingWaitTime        int    `mapstructure:"max_ping_wait_time_ms"`
	V2                     V2     `mapstructure:"v2"`
}

type Errors struct {
	NotEnoughPortsFound    string `mapstructure:"not_enough_ports_found"`
	LoadFromLinkError      string `mapstructure:"load_from_link_error"`
	ConfigNotFound         string `mapstructure:"config_not_found"`
	GroupNotFound          string `mapstructure:"group_not_found"`
	ConfigFormatError      string `mapstructure:"config_format_error"`
	GroupCreatingError     string `mapstructure:"group_creating_error"`
	DirectoryCreatingError string `mapstructure:"directory_creating_error"`
	FileCreatingError      string `mapstructure:"file_creating_error"`
	HttpStatus             string `mapstructure:"http_status"`
	FileRemoveError        string `mapstructure:"file_remove_error"`
	WriteToFileError       string `mapstructure:"write_to_file_error"`
	ConfigRemoveError      string `mapstructure:"config_remove_error"`
	InvalidFormatter       string `mapstructure:"invalid_formatter"`
	InvalidPortSetter      string `mapstructure:"invalid_port_setter"`
	ListenersCountError    string `mapstructure:"listeners_count_error"`
	CollectiveFormatError  string `mapstructure:"collective_format_error"`
	ProxyVariablesSetError string `mapstructure:"proxy_variables_set_error"`
	V2rayActivateError     string `mapstructure:"v2ray_activate_error"`
	V2rayDeactivateError   string `mapstructure:"v2ray_deactivate_error"`
}

type V2 struct {
	ConfigurationPath string `mapstructure:"config_path"`
	DefaultPort       uint16 `mapstructure:"default_port"`
	DefaultHost       string `mapstructure:"default_host"`
	DNS               string `mapstructure:"dns"`
	Inbounds          string `mapstructure:"inbounds"`
	Log               string `mapstructure:"log"`
	Policy            string `mapstructure:"policy"`
	Routing           string `mapstructure:"routing"`
	Stats             string `mapstructure:"stats"`
}

func getDefaultConfig() string {
	return "config.json"
}

func NewConfig() (*Configuration, error) {
	path := os.Getenv("CFG_PATH")
	if path == "" {
		path = getDefaultConfig()
	}
	viperConfig := viper.New()
	viperConfig.SetConfigFile(path)
	viperConfig.AutomaticEnv()
	if err := viperConfig.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("config file not found %w", err)
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
