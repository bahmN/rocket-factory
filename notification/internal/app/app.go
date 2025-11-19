package app

import (
	"context"
	"fmt"

	"github.com/bahmN/rocket-factory/notification/internal/config"
	"github.com/bahmN/rocket-factory/platform/pkg/closer"
	"github.com/bahmN/rocket-factory/platform/pkg/logger"
	"golang.org/x/sync/errgroup"
)

type App struct {
	diContainer *diContainer
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
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		logger.Info(ctx, "Starting order consumer service (OrderPaid)")
		if err := a.runPaidConsumer(gCtx); err != nil {
			return fmt.Errorf("order consumer service error: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		logger.Info(ctx, "Starting order consumer service (ShipAssembled)")
		if err := a.runAssembledConsumer(gCtx); err != nil {
			return fmt.Errorf("order consumer service error: %w", err)
		}
		return nil
	})

	select {
	case <-gCtx.Done():
		return gCtx.Err()
	default:
	}

	return g.Wait()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
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

func (a *App) runPaidConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderPaid Kafka consumer запущен")

	service := a.diContainer.OrderPaidConsumerService()
	err := service.RunConsumer(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) runAssembledConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderAssembled Kafka consumer запущен")

	service := a.diContainer.OrderAssembledConsumerService()
	err := service.RunConsumer(ctx)
	if err != nil {
		return err
	}
	return nil
}
