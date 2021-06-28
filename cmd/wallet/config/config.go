package config

import "github.com/jhdriver/UBChain/config"

type Config struct {
	ConfigFile  string
	Format      bool
	TestNet     bool
	KeyStoreDir string
	config.RpcConfig
}
