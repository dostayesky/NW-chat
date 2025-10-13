package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
)

type Config interface {
	App() App
	RABBITMQ() RABBITMQ
	Database() Database
	JWT() JWT
	String() string
}

type App struct {
	Port string
}

type RABBITMQ struct {
	URI string
}

type Database struct {
	DSN  string
	Name string
}

type JWT struct {
	Token string
}

type config struct {
	AppCfg      App
	DatabaseCfg Database
	RabbitMqCfg RABBITMQ
	JwtCfg      JWT
}

func (c *config) App() App           { return c.AppCfg }
func (c *config) Database() Database { return c.DatabaseCfg }
func (c *config) JWT() JWT           { return c.JwtCfg }
func (c *config) RABBITMQ() RABBITMQ { return c.RabbitMqCfg }

func (c *config) String() string {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatalf("Failed to convert config to JSON: %v", err)
	}
	return string(jsonBytes)
}

func InitConfig() (Config, error) {
	cfg := &config{
		AppCfg: App{
			Port: getEnv("ADDR", ""),
		},
		DatabaseCfg: Database{
			DSN:  getEnv("DATABASE_DSN", ""),
			Name: getEnv("DATABASE_NAME", ""),
		},
		JwtCfg: JWT{
			Token: getEnv("JWT_SECRET", ""),
		},
		RabbitMqCfg: RABBITMQ{
			URI: getEnv("RABBITMQ_URI", ""),
		},
	}

	if err := validator.New().Struct(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
