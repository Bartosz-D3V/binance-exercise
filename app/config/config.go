package config

import "github.com/caarlos0/env/v6"

type Config struct {
	Symbol         string `env:"SYMBOL" envDefault:"btcusdt"`
	QuantityToSell string `env:"QUANTITY_TO_SELL" envDefault:"25.0"`
	MinimumBid     string `env:"MINIMUM_BID" envDefault:"2200.0"`
}

func New() (*Config, error) {
	cfg := &Config{}
	return cfg, env.Parse(cfg)
}
