package config

import (
	"log"
	"os"

	"gopkg.in/gcfg.v1"
)

type MainConfig struct {
	Server struct {
		Name string
		Port string
	}

	DBConfig struct {
		SlaveDSN      string
		MasterDSN     string
		RetryInterval int
		MaxIdleConn   int
		MaxConn       int
	}
}

func ReadModuleConfig(cfg interface{}, path string, module string) bool {
	environ := os.Getenv("TKPENV")
	if environ == "" {
		environ = "development"
	}

	fname := path + "/" + module + "." + environ + ".ini"
	err := gcfg.ReadFileInto(cfg, fname)
	if err == nil {
		return true
	}
	return false
}

func ReadConfig(cfg interface{}, module string) interface{} {
	ok := ReadModuleConfig(cfg, "files/etc/tauria", module)
	if !ok {
		log.Fatalln("failed to read config for ", module)
	}

	return cfg
}
