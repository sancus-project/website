package main

import (
	"time"
)

const (
	DefaultPIDFile         = "/tmp/tableflip.pid"
	DefaultPort            = 8080
	DefaultReadTimeout     = 5 * time.Second
	DefaultWriteTimeout    = 5 * time.Second
	DefaultIdleTimeout     = 30 * time.Second
	DefaultGracefulTimeout = 60 * time.Second
)

type ServerConfig struct {
	Development     bool
	PIDFile         string
	Port            uint16
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	GracefulTimeout time.Duration
}

func NewConfig() ServerConfig {
	return ServerConfig{
		PIDFile:      DefaultPIDFile,
		Port:         DefaultPort,
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
		IdleTimeout:  DefaultIdleTimeout,
	}
}
