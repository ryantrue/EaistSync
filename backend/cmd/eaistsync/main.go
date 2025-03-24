// main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ryantrue/EaistSync/pkg/api/rest"
	"github.com/ryantrue/EaistSync/pkg/config"
	"github.com/ryantrue/EaistSync/pkg/cron"
	"github.com/ryantrue/EaistSync/pkg/db"
	"github.com/ryantrue/EaistSync/pkg/logger"
	"github.com/ryantrue/EaistSync/pkg/messaging"
	"github.com/ryantrue/EaistSync/pkg/migrate"
	"github.com/ryantrue/EaistSync/pkg/server"
	"github.com/ryantrue/EaistSync/pkg/telegrambot"
	"github.com/ryantrue/EaistSync/pkg/utils"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// processedContractIDs хранит идентификаторы уже обработанных контрактов.
var processedContractIDs = make(map[int64]bool)

func main() {
	// Создаем корневой контекст для graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка выполнения приложения: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Инициализация логгера.
	log, err := logger.NewLogger()
	if err != nil {
		return fmt.Errorf("не удалось инициализировать логгер: %w", err)
	}
	defer log.Sync()

	// Загрузка конфигурации.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации", zap.Error(err))
	}

	// Подключаемся к PostgreSQL.
	dbConn, err := sqlx.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		log.Fatal("Ошибка открытия БД", zap.Error(err))
	}
	defer dbConn.Close()

	// Запуск миграций.
	migrator := migrate.NewMigrator(dbConn, "./migrations", log)
	if err := migrator.RunUp(); err != nil {
		log.Fatal("Ошибка запуска миграций", zap.Error(err))
	}

	// Создаем HTTP-клиент для REST API.
	httpClient, err := rest.NewHTTPClient(30 * time.Second)
	if err != nil {
		log.Fatal("Ошибка создания HTTP клиента", zap.Error(err))
	}

	// Инициализируем Kafka продюсера с повторными попытками.
	producer, err := initKafkaProducer(ctx, []string{cfg.KafkaBrokers}, "eaist_updates", log)
	if err != nil {
		log.Fatal("Ошибка создания Kafka продюсера", zap.Error(err))
	}
	defer producer.Close()

	// Инициализируем Telegram-бота, если задан токен.
	var telegramBot *telegrambot.TelegramBot
	if cfg.TelegramBotToken != "" {
		telegramBot, err = telegrambot.NewTelegramBot(cfg.TelegramBotToken, cfg.TelegramChatID, 3, 2*time.Second)
		if err != nil {
			log.Error("Ошибка создания Telegram-бота", zap.Error(err))
		}
	}

	// Первичное обновление данных.
	log.Info("Первичный запуск обновления данных")
	start := time.Now()
	newContracts, err := updateData(ctx, httpClient, dbConn, log, producer, cfg)
	if err != nil {
		log.Error("Ошибка при первоначальном обновлении данных", zap.Error(err))
		if telegramBot != nil {
			telegramBot.Notify(ctx, fmt.Sprintf("Первичное обновление данных завершено с ошибкой.\nВремя выполнения: %v\nОшибка: %v", time.Since(start), err))
		}
	} else {
		if len(newContracts) > 0 && telegramBot != nil {
			if err := telegramBot.SendJSONDocument(ctx, newContracts); err != nil {
				log.Error("Ошибка отправки новых контрактов через Telegram", zap.Error(err))
			}
		}
		log.Info("Первичное обновление данных прошло успешно")
	}

	// Создаем новый Scheduler из пакета cron с передачей корневого контекста.
	scheduler := cron.NewScheduler(ctx, log)
	_, err = scheduler.AddTask("@daily", dataUpdater(ctx, httpClient, dbConn, log, producer, cfg, telegramBot))
	if err != nil {
		log.Fatal("Ошибка добавления cron-задачи", zap.Error(err))
	}
	go scheduler.Start()

	// Запуск HTTP-сервера.
	serverAddr := fmt.Sprintf(":%s", cfg.Port)
	appServer := server.NewServer(dbConn, log, cfg)
	serverErrCh := make(chan error, 1)
	go func() {
		log.Info("Запуск сервера", zap.String("addr", serverAddr))
		if err := appServer.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			serverErrCh <- err
		}
	}()

	// Ожидание сигнала завершения.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	select {
	case <-sigCh:
		log.Info("Получен сигнал завершения работы")
	case err := <-serverErrCh:
		log.Error("Ошибка работы сервера", zap.Error(err))
	}

	return nil
}

// initKafkaProducer пытается создать Kafka продюсера с несколькими попытками.
func initKafkaProducer(ctx context.Context, brokers []string, topic string, log *zap.Logger) (messaging.KafkaProducerInterface, error) {
	var producer messaging.KafkaProducerInterface
	var err error
	maxAttempts := 3
	for i := 1; i <= maxAttempts; i++ {
		producer, err = messaging.NewKafkaProducer(brokers, topic, log)
		if err == nil {
			log.Info("Kafka продюсер успешно создан", zap.Int("attempt", i))
			return producer, nil
		}
		log.Warn("Не удалось создать Kafka продюсера", zap.Int("attempt", i), zap.Error(err))
		time.Sleep(5 * time.Second)
	}
	return nil, fmt.Errorf("не удалось создать Kafka продюсера после %d попыток: %w", maxAttempts, err)
}

// updateData обновляет данные из EAIST REST API, сохраняет их в БД, публикует событие в Kafka
// и возвращает список новых контрактов (тех, чьи ID ранее не были обработаны).
func updateData(ctx context.Context, client *http.Client, dbConn *sqlx.DB, log *zap.Logger, producer messaging.KafkaProducerInterface, cfg *config.Config) ([]map[string]interface{}, error) {
	// Авторизация через REST API.
	if err := rest.Login(ctx, client, cfg); err != nil {
		return nil, fmt.Errorf("ошибка авторизации: %w", err)
	}

	// Получение контрактов.
	contracts, err := rest.FetchAllContracts(ctx, client, log, cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения контрактов: %w", err)
	}

	// Получение состояний.
	states, err := rest.FetchStates(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения состояний: %w", err)
	}

	// Фильтруем новые контракты (те, чей ID ранее не встречался).
	var newContracts []map[string]interface{}
	for _, contract := range contracts {
		id, err := utils.ExtractID(contract)
		if err != nil {
			log.Warn("Невозможно извлечь ID контракта", zap.Error(err), zap.Any("contract", contract))
			continue
		}
		if _, exists := processedContractIDs[id]; !exists {
			newContracts = append(newContracts, contract)
		}
	}
	// Обновляем список обработанных контрактов.
	for _, contract := range contracts {
		id, err := utils.ExtractID(contract)
		if err == nil {
			processedContractIDs[id] = true
		}
	}

	// Создаем JSONUpserter с динамическим списком разрешённых таблиц.
	upserter := db.NewJSONUpserter(dbConn, log, []string{"contracts", "states"})

	// Сохраняем данные в БД через новый интерфейс.
	if err := upserter.UpsertMany(ctx, "contracts", contracts); err != nil {
		return nil, fmt.Errorf("ошибка сохранения контрактов: %w", err)
	}
	if err := upserter.UpsertMany(ctx, "states", states); err != nil {
		return nil, fmt.Errorf("ошибка сохранения состояний: %w", err)
	}

	// Формируем сообщение для Kafka.
	updateMessage := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC1123),
		"contracts": len(contracts),
		"states":    len(states),
		"event":     "data_updated",
	}
	if err := producer.PublishMessage(ctx, updateMessage); err != nil {
		return nil, fmt.Errorf("ошибка отправки сообщения в Kafka: %w", err)
	}

	return newContracts, nil
}

// dataUpdater возвращает функцию UpdaterFunc, которая замыкает все необходимые зависимости.
func dataUpdater(ctx context.Context, client *http.Client, dbConn *sqlx.DB, log *zap.Logger, producer messaging.KafkaProducerInterface, cfg *config.Config, telegramBot *telegrambot.TelegramBot) cron.UpdaterFunc {
	return func(ctx context.Context) {
		log.Info("Запуск обновления данных из EAIST")
		newContracts, err := updateData(ctx, client, dbConn, log, producer, cfg)
		if err != nil {
			log.Error("Ошибка обновления данных", zap.Error(err))
			if telegramBot != nil {
				telegramBot.Notify(ctx, fmt.Sprintf("Ошибка обновления данных: %v", err))
			}
		} else if telegramBot != nil {
			if len(newContracts) > 0 {
				if err := telegramBot.SendJSONDocument(ctx, newContracts); err != nil {
					log.Error("Ошибка отправки новых контрактов через Telegram", zap.Error(err))
				}
			} else {
				log.Info("Обновление данных выполнено, новых контрактов не обнаружено")
			}
		}
	}
}
