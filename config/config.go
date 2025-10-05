package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		ServerHost string `mapstructure:"server_host"`
		HttpPort   int    `mapstructure:"http_port"`
		CORS       struct {
			AllowOrigins        []string      `mapstructure:"allow_origins"`
			AllowMethods        []string      `mapstructure:"allow_methods"`
			AllowHeaders        []string      `mapstructure:"allow_headers"`
			AllowCredentials    bool          `mapstructure:"allow_credentials"`
			AllowWebSockets     bool          `mapstructure:"allow_websockets"`
			AllowFiles          bool          `mapstructure:"allow_files"`
			AllowPrivateNetwork bool          `mapstructure:"allow_private_network"`
			MaxAge              time.Duration `mapstructure:"max_age"`
			ExposeHeaders       []string      `mapstructure:"expose_headers"`
		} `mapstructure:"cors"`
		Http struct {
			WriteTimeout   time.Duration `mapstructure:"write_timeout"`
			ReadTimeout    time.Duration `mapstructure:"read_timeout"`
			IdleTimeout    time.Duration `mapstructure:"idle_timeout"`
			MaxHeaderBytes int           `mapstructure:"max_header_bytes"`
		} `mapstructure:"http"`
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

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
