package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/bahmN/rocket-factory/order/internal/config"
	"github.com/bahmN/rocket-factory/platform/pkg/closer"
	"github.com/bahmN/rocket-factory/platform/pkg/logger"
	"github.com/bahmN/rocket-factory/platform/pkg/migrator"
	orderV1 "github.com/bahmN/rocket-factory/shared/pkg/openapi/order/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
	listener    net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Запускаем консьюмер
	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("consumer error: %w", err)
		}
	}()

	// Запускаем HTTP-сервер
	go func() {
		if err := a.runHTTPServer(ctx); err != nil {
			errCh <- fmt.Errorf("http server error: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		logger.Error(ctx, "❌ Компонент завершился с ошибкой, завершение работы", zap.Error(err))
		// Триггерим cancel, чтобы остановить второй компонент
		cancel()
		// Дождись завершения всех задач (если есть graceful shutdown внутри)
		<-ctx.Done()
		return err
	case <-ctx.Done():
		logger.Info(ctx, "🔔 Получен сигнал завершения работы")
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initHTTPServer,
		a.initMigrations,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJSON(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().OrderHTTP.Address())
	if err != nil {
		return err
	}
	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		ler := listener.Close()
		if ler != nil && !errors.Is(ler, net.ErrClosed) {
			return ler
		}

		return nil
	})

	a.listener = listener

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	server, err := orderV1.NewServer(a.diContainer.InventoryV1API(ctx))
	if err != nil {
		return err
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		err = a.httpServer.Shutdown(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", server)

	a.httpServer = &http.Server{
		Addr:              config.AppConfig().OrderHTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("HTTP OrderService server listening on %s", config.AppConfig().OrderHTTP.Address()))

	err := a.httpServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initMigrations(ctx context.Context) error {
	pool := a.diContainer.InitDBPool(ctx)
	dir := config.AppConfig().Postgres.MigrationsDir()
	runner := migrator.NewMigrator(stdlib.OpenDB(*pool.Config().ConnConfig), dir)

	err := runner.Up()
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("db migration error: %v", err))
		return nil
	}

	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderAssembled Kafka consumer запущен")

	err := a.diContainer.OrderConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
