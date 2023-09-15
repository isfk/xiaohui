package config

import (
	"bytes"
	_ "embed"
	"log/slog"

	"github.com/spf13/viper"
)

var Config config

type config struct {
	XiaoHui XiaoHuiConf
	Cron    struct {
		Name string    `json:"name"`
		Jobs []JobConf `json:"jobs"`
		Log  LogConf
	} `json:"cron"`
}

type XiaoHuiConf struct {
	Path      string `json:"path"`
	AudioPath string `json:"audio_path"`
	DocsPath  string `json:"docs_path"`
}

type JobConf struct {
	Name    string `json:"name"`
	Spec    string `json:"spec"`
	Enabled bool   `json:"enabled"`
}

type LogConf struct {
	DefaultPath string `json:"default_path"`
	ErrorPath   string `json:"error_path"`
	MaxSize     int    `json:"max_size"`
}

type Conf struct {
	Config config
}

//go:embed config.yaml
var f []byte

func NewConfig() *Conf {
	var config config
	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(bytes.NewReader(f)); err != nil {
		slog.Error("config %v", slog.Any("err", err.Error()))
		return &Conf{}
	}

	if err := viper.Unmarshal(&config); err != nil {
		slog.Error("config %v", slog.Any("err", err.Error()))
		return &Conf{}
	}

	slog.Info("config %v", slog.Any("配置", config))

	return &Conf{Config: config}
}
