// pkg/cron/scheduler.go
package cron

import (
	"context"
	"runtime/debug"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// zapCronLogger адаптирует zap-события к интерфейсу cron.Logger.
type zapCronLogger struct {
	sugar *zap.SugaredLogger
}

// Info реализует метод Info интерфейса cron.Logger.
func (z *zapCronLogger) Info(msg string, keysAndValues ...interface{}) {
	z.sugar.Infow(msg, keysAndValues...)
}

// Error реализует метод Error интерфейса cron.Logger.
func (z *zapCronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	// Добавляем ошибку в набор ключ-значение.
	z.sugar.Errorw(msg, append([]interface{}{"error", err}, keysAndValues...)...)
}

// UpdaterFunc — функция, которую будет выполнять cron-задача.
type UpdaterFunc func(ctx context.Context)

// Scheduler инкапсулирует cron-планировщик и использует базовый контекст для всех задач.
type Scheduler struct {
	cron   *cron.Cron
	logger *zap.Logger
	ctx    context.Context
}

// NewScheduler создаёт новый Scheduler с настройками по умолчанию и базовым контекстом.
// Логирование интегрировано через cron.WithLogger, что позволяет использовать zap для подробного вывода.
func NewScheduler(ctx context.Context, logger *zap.Logger) *Scheduler {
	c := cron.New(
		cron.WithChain(
			cron.Recover(cron.DefaultLogger), // базовое восстановление после паники
		),
		cron.WithLogger(&zapCronLogger{sugar: logger.Sugar()}),
	)
	return &Scheduler{
		cron:   c,
		logger: logger,
		ctx:    ctx,
	}
}

// AddTask добавляет новую задачу с заданным расписанием и обработчиком.
// Возвращает идентификатор задачи или ошибку.
func (s *Scheduler) AddTask(schedule string, updater UpdaterFunc) (cron.EntryID, error) {
	// Оборачиваем выполнение задачи для логирования и обработки паники.
	task := func() {
		// Используем базовый контекст, чтобы задачи могли реагировать на отмену.
		ctx := s.ctx
		s.logger.Info("Cron job started", zap.String("schedule", schedule))
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error("Cron job panicked",
					zap.Any("recover", r),
					zap.String("schedule", schedule),
					zap.ByteString("stack", debug.Stack()),
				)
			} else {
				s.logger.Info("Cron job finished", zap.String("schedule", schedule))
			}
		}()
		updater(ctx)
	}
	id, err := s.cron.AddFunc(schedule, task)
	if err != nil {
		s.logger.Error("Error adding cron task", zap.Error(err), zap.String("schedule", schedule))
		return 0, err
	}
	return id, nil
}

// Start запускает планировщик и блокирует выполнение до отмены базового контекста.
// При отмене контекста планировщик корректно останавливается.
func (s *Scheduler) Start() {
	s.logger.Info("Starting cron scheduler")
	s.cron.Start()
	<-s.ctx.Done()
	s.logger.Info("Stopping cron scheduler")
	s.cron.Stop()
}
