package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Container struct {
		App *App
		// Redis *Redis
		DB   *DB
		HTTP *HTTP
		// RabbitMQ *RabbitMQ
	}
	// App contains all the environment variables for the application
	App struct {
		Name string
		Mode string
	}

	// Redis contains all the environment variables for the cache service
	Redis struct {
		Addr     string
		Password string
	}
	// DB contains all the environment variables for the database
	DB struct {
		Connection string
		Host       string
		Port       string
		User       string
		Password   string
		Name       string
	}
	// HTTP contains all the environment variables for the http server
	HTTP struct {
		Env            string
		URL            string
		Port           string
		AllowedOrigins string
		SlipAPIURL     string
	}

	RabbitMQ struct {
		AmqpProtocal string
		AmqpAddr     string
		AmqpUser     string
		AmqpPassword string
		AmqpPort     string
		AmqpVhost    string

		MqttProtocal string
		MqttAddr     string
		MqttUser     string
		MqttPassword string
		MqttPort     string
		MqttVhost    string
	}
)

// New loads the core application configuration (DB, HTTP, etc.)
func New() (*Container, error) {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			return nil, err
		}
	}

	app := &App{
		Name: os.Getenv("APP_NAME"),
		Mode: os.Getenv("APP_MODE"),
	}

	db := &DB{
		Connection: os.Getenv("DB_CONNECTION"),
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		Name:       os.Getenv("DB_NAME"),
	}

	http := &HTTP{
		Env:            os.Getenv("APP_ENV"),
		URL:            os.Getenv("HTTP_URL"),
		Port:           os.Getenv("HTTP_PORT"),
		AllowedOrigins: os.Getenv("HTTP_ALLOWED_ORIGINS"),
		SlipAPIURL:     os.Getenv("SLIP_API_URL"),
	}

	return &Container{
		App:  app,
		DB:   db,
		HTTP: http,
	}, nil
}


