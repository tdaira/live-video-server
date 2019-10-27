package config

import (
	"github.com/spf13/viper"
	"log"
)

// 環境変数から取得するデータの構造体
type Config struct {
	LogLevel string `mapstructure:"log_level"`
	Server   Server `mapstructure:"server"`
	Source   Source `mapstructure:"source"`
}

// gRPCサーバに関する設定
type Server struct {
	Port           string `mapstructure:"port"`
	ConnectTimeout uint   `mapstructure:"connect_timeout"`
	KeepAlive      uint   `mapstructure:"keep_alive"`
	MetricsPort    string `mapstructure:"metrics_port"`
}

// 入力データに関する設定
type Source struct {
	VideoPath   string `mapstructure:"video_path"`
}

// Tomlファイルから設定を取得
func SetupConfig() (Config, error) {
	// Set config path.
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Couldn't read config: %s", err)
	}

	return cfg, nil
}
