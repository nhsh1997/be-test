//AUTO-GENERATED: DO NOT EDIT

package config

import (
	"bytes"
	"strings"

	"git.begroup.team/platform-core/kitchen/l"
	"github.com/spf13/viper"
)

var (
	ll = l.New()
)

type Tracing struct {
	Enabled bool
}

type Base struct {
	HTTPAddress int `yaml:"http_address" mapstructure:"http_address"`
	GRPCAddress int `mapstructure:"grpc_address"`
	Environment string

	Tracing Tracing
}

func Load() *Config {
	var cfg = &Config{}

	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	if err != nil {
		ll.Fatal("Failed to read viper config", l.Error(err))
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	err = viper.Unmarshal(&cfg)
	if err != nil {
		ll.Fatal("Failed to unmarshal config", l.Error(err))
	}

	ll.Info("Config loaded", l.Object("config", cfg))
	return cfg
}
