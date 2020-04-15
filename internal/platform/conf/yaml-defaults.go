package conf

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

var defaults = []byte(`
web:
  server:
    readTimeout: 3s
    gracefulShutdownTimeout: 1s

  client:
    readTimeout: 3s
    writeTimeout: 1s
    maxConnectionsDuration: 10m
    maxIdleConnectionsPerHost: 512
    maxIdleConnectionDuration: 10s

logging:
  enable: true
  timestamp: true
  level: info
  format: json

cache:
  mappings:
    maxSize: 100
  query:
    maxSize: 100
  parser:
    maxSize: 100

database:
  timeout: 1000
`)

func readDefaults(cfg *Config) {
	err := yaml.Unmarshal(defaults, cfg)
	if err != nil {
		log.Printf("failed to unmarshal defaults : %s", err)
	}
}
