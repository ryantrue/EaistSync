package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Username       string
	Password       string
	APIType        string // "rest"
	DatabaseDSN    string
	Port           string
	KafkaBrokers   string
	MinioEndpoint  string // Параметры для MinIO
	MinioAccessKey string // Параметры для MinIO
	MinioSecretKey string // Параметры для MinIO

	// Параметры для REST API
	ContractsURL   string
	PageSize       int
	MaxConcurrency int
	LoginURL       string

	// Параметры для Telegram-бота
	TelegramBotToken string
	TelegramChatID   int64

	// JWT секрет для подписи токенов
	JWTSecret string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = viper.ReadInConfig()

	username := viper.GetString("USERNAME")
	password := viper.GetString("PASSWORD")
	apiType := viper.GetString("API_TYPE")
	dbdsn := viper.GetString("DATABASE_DSN")
	port := viper.GetString("PORT")
	kafkaBrokers := viper.GetString("KAFKA_BROKERS")

	// Чтение параметров для MinIO
	minioEndpoint := viper.GetString("MINIO_ENDPOINT")
	minioAccessKey := viper.GetString("MINIO_ROOT_USER")
	minioSecretKey := viper.GetString("MINIO_ROOT_PASSWORD")

	// Чтение параметров для REST API
	contractsURL := viper.GetString("CONTRACTS_URL")
	pageSize := viper.GetInt("PAGE_SIZE")
	maxConcurrency := viper.GetInt("MAX_CONCURRENCY")
	loginURL := viper.GetString("LOGIN_URL")

	// Чтение параметров для Telegram-бота
	telegramBotToken := viper.GetString("TELEGRAM_BOT_TOKEN")
	telegramChatID := int64(viper.GetInt("TELEGRAM_CHAT_ID"))

	// Чтение JWT секрета
	jwtSecret := viper.GetString("JWT_SECRET")

	// Проверка и установка значений по умолчанию
	if username == "" || password == "" {
		return nil, fmt.Errorf("USERNAME или PASSWORD не заданы")
	}
	if apiType == "" {
		apiType = "rest"
	}
	if dbdsn == "" {
		return nil, fmt.Errorf("DATABASE_DSN не задан")
	}
	if port == "" {
		port = "8080"
	}
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}
	if minioEndpoint == "" {
		minioEndpoint = "minio:9000"
	}
	if minioAccessKey == "" {
		minioAccessKey = "minioadmin"
	}
	if minioSecretKey == "" {
		minioSecretKey = "minioadmin"
	}
	if contractsURL == "" {
		contractsURL = "https://eaist.mos.ru/eaist2rc/api/contracts/contract/list"
	}
	if pageSize == 0 {
		pageSize = 500
	}
	if maxConcurrency == 0 {
		maxConcurrency = 5
	}
	if loginURL == "" {
		loginURL = "https://eaist.mos.ru/module/protected-admin/api/login"
	}
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET не задан")
	}

	// Параметры для Telegram-бота можно не задавать, если функциональность не требуется.
	// Но если задан токен, то идентификатор чата должен быть установлен.
	if telegramBotToken != "" && telegramChatID == 0 {
		return nil, fmt.Errorf("TELEGRAM_CHAT_ID не задан для Telegram-бота")
	}

	return &Config{
		Username:         username,
		Password:         password,
		APIType:          apiType,
		DatabaseDSN:      dbdsn,
		Port:             port,
		KafkaBrokers:     kafkaBrokers,
		MinioEndpoint:    minioEndpoint,
		MinioAccessKey:   minioAccessKey,
		MinioSecretKey:   minioSecretKey,
		ContractsURL:     contractsURL,
		PageSize:         pageSize,
		MaxConcurrency:   maxConcurrency,
		LoginURL:         loginURL,
		TelegramBotToken: telegramBotToken,
		TelegramChatID:   telegramChatID,
		JWTSecret:        jwtSecret,
	}, nil
}
