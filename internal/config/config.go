package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type BalancerConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":80"`
	RedisAddress  string `env:"REDIS_ADDRESS" envDefault:":6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`
}

var balancerConf *BalancerConfig

func GetConfig() *BalancerConfig {
	if balancerConf == nil {
		balancerConf = &BalancerConfig{}
		if err := env.Parse(balancerConf); err != nil {
			return nil
		}
		flag.StringVar(&balancerConf.ServerAddress, "a", balancerConf.ServerAddress, "address:port to listen")
		flag.StringVar(&balancerConf.RedisAddress, "d", balancerConf.RedisAddress, "address:port of Redis")
		flag.StringVar(&balancerConf.RedisPassword, "p", balancerConf.RedisPassword, "Redis password")
		flag.Parse()
	}

	return balancerConf
}
