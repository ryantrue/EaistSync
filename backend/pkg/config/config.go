package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/viper"
)

// Config хранит конфигурационные параметры приложения.
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

// getValue пытается получить значение из переменной окружения.
// Если значение не найдено, проверяет переменную с суффиксом _FILE и считывает содержимое файла.
func getValue(key string) (string, error) {
	value := viper.GetString(key)
	if value != "" {
		return value, nil
	}

	// Пытаемся получить путь к файлу-секрету
	filePath := viper.GetString(key + "_FILE")
	if filePath == "" {
		return "", nil // значение не задано ни напрямую, ни через файл
	}

	// Читаем содержимое файла
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("не удалось прочитать файл %s: %w", filePath, err)
	}
	return strings.TrimSpace(string(data)), nil
}

// LoadConfig загружает конфигурацию из файла .env и переменных окружения.
func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = viper.ReadInConfig()

	username, err := getValue("USERNAME")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении USERNAME: %w", err)
	}
	password, err := getValue("PASSWORD")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении PASSWORD: %w", err)
	}
	apiType, err := getValue("API_TYPE")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении API_TYPE: %w", err)
	}
	dbdsn, err := getValue("DATABASE_DSN")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении DATABASE_DSN: %w", err)
	}
	port, err := getValue("PORT")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении PORT: %w", err)
	}
	kafkaBrokers, err := getValue("KAFKA_BROKERS")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении KAFKA_BROKERS: %w", err)
	}

	// Чтение параметров для MinIO
	minioEndpoint, err := getValue("MINIO_ENDPOINT")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении MINIO_ENDPOINT: %w", err)
	}
	minioAccessKey, err := getValue("MINIO_ROOT_USER")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении MINIO_ROOT_USER: %w", err)
	}
	minioSecretKey, err := getValue("MINIO_ROOT_PASSWORD")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении MINIO_ROOT_PASSWORD: %w", err)
	}

	// Чтение параметров для REST API
	contractsURL, err := getValue("CONTRACTS_URL")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении CONTRACTS_URL: %w", err)
	}
	// Для числовых значений используем viper напрямую (если они заданы в .env)
	pageSize := viper.GetInt("PAGE_SIZE")
	maxConcurrency := viper.GetInt("MAX_CONCURRENCY")
	loginURL, err := getValue("LOGIN_URL")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении LOGIN_URL: %w", err)
	}

	// Чтение параметров для Telegram-бота с использованием getValue.
	telegramBotToken, err := getValue("TELEGRAM_BOT_TOKEN")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении TELEGRAM_BOT_TOKEN: %w", err)
	}
	telegramChatIDStr, err := getValue("TELEGRAM_CHAT_ID")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении TELEGRAM_CHAT_ID: %w", err)
	}
	// Устанавливаем значение в viper и получаем его как int64
	viper.Set("TELEGRAM_CHAT_ID", telegramChatIDStr)
	telegramChatID := viper.GetInt64("TELEGRAM_CHAT_ID")

	// Если задан Telegram Bot Token, но TELEGRAM_CHAT_ID не установлен,
	// отключаем функциональность Telegram, чтобы не прерывать работу приложения.
	if telegramBotToken != "" && telegramChatID == 0 {
		telegramBotToken = ""
		telegramChatID = 0
	}

	// Чтение JWT секрета
	jwtSecret, err := getValue("JWT_SECRET")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении JWT_SECRET: %w", err)
	}

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
