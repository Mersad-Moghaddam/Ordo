package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type ApplicationConfiguration struct {
	ServiceName      string
	HTTPPort         int
	MySQLDSN         string
	RedisAddress     string
	GracefulShutdown time.Duration
}

type Option func(applicationConfiguration *ApplicationConfiguration) error

func NewApplicationConfiguration(options ...Option) (ApplicationConfiguration, error) {
	applicationConfiguration := ApplicationConfiguration{
		ServiceName:      "ordo-api",
		HTTPPort:         8080,
		MySQLDSN:         "root:root@tcp(localhost:3306)/ordo?parseTime=true",
		RedisAddress:     "localhost:6379",
		GracefulShutdown: 10 * time.Second,
	}

	for _, option := range options {
		if option == nil {
			continue
		}
		if optionError := option(&applicationConfiguration); optionError != nil {
			return ApplicationConfiguration{}, optionError
		}
	}

	return applicationConfiguration, nil
}

func WithHTTPPort(httpPort int) Option {
	return func(applicationConfiguration *ApplicationConfiguration) error {
		if httpPort <= 0 {
			return fmt.Errorf("http port must be positive")
		}
		applicationConfiguration.HTTPPort = httpPort
		return nil
	}
}

func WithServiceName(serviceName string) Option {
	return func(applicationConfiguration *ApplicationConfiguration) error {
		if serviceName == "" {
			return fmt.Errorf("service name must not be empty")
		}
		applicationConfiguration.ServiceName = serviceName
		return nil
	}
}

func WithEnvironment() Option {
	return func(applicationConfiguration *ApplicationConfiguration) error {
		httpPortValue := os.Getenv("ORDO_HTTP_PORT")
		if httpPortValue != "" {
			httpPort, conversionError := strconv.Atoi(httpPortValue)
			if conversionError != nil {
				return fmt.Errorf("invalid ORDO_HTTP_PORT: %w", conversionError)
			}
			if httpPort <= 0 {
				return fmt.Errorf("ORDO_HTTP_PORT must be positive")
			}
			applicationConfiguration.HTTPPort = httpPort
		}

		mysqlDSNValue := os.Getenv("ORDO_MYSQL_DSN")
		if mysqlDSNValue != "" {
			applicationConfiguration.MySQLDSN = mysqlDSNValue
		}

		redisAddressValue := os.Getenv("ORDO_REDIS_ADDRESS")
		if redisAddressValue != "" {
			applicationConfiguration.RedisAddress = redisAddressValue
		}

		return nil
	}
}
