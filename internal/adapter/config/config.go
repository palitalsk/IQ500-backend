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
		AI   *AI
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
	// Database contains all the environment variables for the database
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
	}

	// AI contains all the environment variables for the AI service
	AI struct {
		BaseURL string
		Timeout int
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

func New() (*Container, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	app := &App{
		Name: os.Getenv("APP_NAME"),
		Mode: os.Getenv("APP_MODE"),
	}

	// redis := &Redis{
	// 	Addr:     os.Getenv("REDIS_ADDR"),
	// 	Password: os.Getenv("REDIS_PASSWORD"),
	// }

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
	}

	ai := &AI{
		BaseURL: os.Getenv("AI_BASE_URL"),
	}

	// rabbitMQ := &RabbitMQ{
	// 	AmqpProtocal: os.Getenv("RABBITMQ_AMQP_PROTOCOL"),
	// 	AmqpAddr:     os.Getenv("RABBITMQ_AMQP_ADDR"),
	// 	AmqpUser:     os.Getenv("RABBITMQ_AMQP_USER"),
	// 	AmqpPassword: os.Getenv("RABBITMQ_AMQP_PASSWORD"),
	// 	AmqpPort:     os.Getenv("RABBITMQ_AMQP_PORT"),
	// 	AmqpVhost:    os.Getenv("RABBITMQ_AMQP_VHOST"),

	// 	MqttProtocal: os.Getenv("RABBITMQ_MQTT_PROTOCOL"),
	// 	MqttAddr:     os.Getenv("RABBITMQ_MQTT_ADDR"),
	// 	MqttUser:     os.Getenv("RABBITMQ_MQTT_USER"),
	// 	MqttPassword: os.Getenv("RABBITMQ_MQTT_PASSWORD"),
	// 	MqttPort:     os.Getenv("RABBITMQ_MQTT_PORT"),
	// 	MqttVhost:    os.Getenv("RABBITMQ_MQTT_VHOST"),
	// }

	return &Container{
		app,
		// redis,
		db,
		http,
		ai,
		// rabbitMQ,
	}, nil
}
