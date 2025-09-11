package config

import "github.com/spf13/viper"

type AppConfig struct {
	App struct {
		ServerHost string `mapstructure:"server_host"`
		HttpPort   int    `mapstructure:"http_port"`
	} `mapstructure:"app"`

	Jwt struct {
		SecretKey   string `mapstructure:"secret_key"`
		AccessName  string `mapstructure:"access_name"`
		RefreshName string `mapstructure:"refresh_name"`
	} `mapstructure:"jwt"`

	Services struct {
		AuthPort    int `mapstructure:"auth_port"`
		UserPort    int `mapstructure:"user_port"`
		ProductPort int `mapstructure:"product_port"`
		PostPort    int `mapstructure:"post_port"`
		ChatPort    int `mapstructure:"chat_port"`
	} `mapstructure:"services"`

	MessageQueue struct {
		RHost     string `mapstructure:"rb_host"`
		RUser     string `mapstructure:"rb_user"`
		RPassword string `mapstructure:"rb_password"`
		RVhost    string `mapstructure:"rb_vhost"`
	} `mapstructure:"message_queue"`
}

func LoadConfig() (*AppConfig, error) {
	viper.SetConfigFile("config/config.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
