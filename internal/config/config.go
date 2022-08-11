package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var cfg *Config

type Config struct {
	Project  Project  `yaml:"project"`
	DB       Database `yaml:"database"`
	Rest     Rest     `yaml:"rest"`
	GRPC     GRPC     `yaml:"gRPC"`
	Logger   Logger   `yaml:"logger"`
	Consumer Consumer `yaml:"consumer"`
}
type Project struct {
	Name        string `yaml:"name"`
	Debug       bool   `yaml:"debug"`
	Environment string `yaml:"environment"`
	Version     string `yaml:"version"`
}

type Logger struct {
	DynamicLevel *zap.AtomicLevel
	InitialLevel string `yaml:"level"`
}

type Database struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	DbName   string `yaml:"db_name"`
}

type Rest struct {
	Host            string `yaml:"host"`
	BusinessPort    string `yaml:"port"`
	DebugPort       string `yaml:"debug_port"`
	DocPort         string `yaml:"documentation_port"`
	GracefulTimeout uint   `yaml:"graceful_timeout"`
}

type GRPC struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	ConnTimeout uint   `yaml:"connection_timeout"`
}

type Consumer struct {
	Type        string   `yaml:"type"`
	Brokers     []string `yaml:"brokers"`
	Topics      []string `yaml:"topics"`
	ConsGroupID string   `yaml:"consumer_group_id"`
}

func GetConfig() Config {
	if cfg != nil {
		return *cfg
	}

	return Config{}
}

func ReadConfigYML(pathToConfFile string) error {
	if cfg != nil {
		return fmt.Errorf("config instance is already exist")
	}

	config := new(Config)
	if err := cleanenv.ReadConfig(pathToConfFile, config); err != nil {
		return fmt.Errorf("unable read config: %w", err)
	}

	loggerLevel, err := zapcore.ParseLevel(config.Logger.InitialLevel)
	if err != nil {
		return fmt.Errorf("wrong logger level: %w", err)
	}
	DynamicLevel := zap.NewAtomicLevelAt(loggerLevel)
	config.Logger.DynamicLevel = &DynamicLevel

	cfg = config

	return nil
}

func SetDynamicLogLevel(level zapcore.Level) {
	if cfg != nil {
		cfg.Logger.DynamicLevel.SetLevel(level)
	}
}
