package config

import "github.com/UBChainNet/UBChain/config"

type Config struct {
	ConfigFile  string
	Format      bool
	TestNet     bool
	KeyStoreDir string
	config.RpcConfig
}
