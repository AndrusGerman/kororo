package services

import (
	"fmt"
	"kororo/internal/adapters/config"
	"kororo/internal/core/ports"
)

type LogService struct {
	conf *config.Config
}

func NewLogService(conf *config.Config) ports.LogService {
	return &LogService{
		conf: conf,
	}
}

func (s *LogService) Info(module string, message string) {
	if s.conf.Debug() {
		fmt.Printf("%s: %s\n", module, message)
	}
}

func (s *LogService) Error(module string, message string) {
	if s.conf.Debug() {
		fmt.Printf("%s: %s\n", module, message)
	}
}
