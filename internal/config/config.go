package config

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

var cfg *Config

type Config struct {
	Project Project  `yaml:"project"`
	DB      Database `yaml:"database"`
	Rest    Rest     `yaml:"rest"`
	GRPC    GRPC     `yaml:"gRPC"`
	Logger  Logger   `yaml:"logger"`
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
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"db_name"`
}

type Rest struct {
	Host            string `yaml:"host"`
	BusinessPort    string `yaml:"port"`
	DebugPort       string `yaml:"debugPort"`
	GracefulTimeout uint   `yaml:"gracefulTimeout"`
}

type GRPC struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
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

	confFile, err := os.Open(filepath.Clean(pathToConfFile))
	if err != nil {
		return fmt.Errorf("can't open config file: %w", err)
	}
	defer func() { _ = confFile.Close() }()

	if err := yaml.NewDecoder(confFile).Decode(&cfg); err != nil {
		return fmt.Errorf("can't decode config: %w", err)
	}

	loggerLevel, err := zapcore.ParseLevel(cfg.Logger.InitialLevel)
	if err != nil {
		return fmt.Errorf("wrong logger level: %w", err)
	}
	DynamicLevel := zap.NewAtomicLevelAt(loggerLevel)
	cfg.Logger.DynamicLevel = &DynamicLevel

	return nil
}

func SetDynamicLogLevel(level zapcore.Level) {
	if cfg != nil {
		cfg.Logger.DynamicLevel.SetLevel(level)
	}
}
