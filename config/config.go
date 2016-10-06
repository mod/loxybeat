// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Address string        `config:"address"`
	Timeout time.Duration `config:"timeout"`
}

var DefaultConfig = Config{
	Address: ":8080",
	Timeout: 20 * time.Second,
}
