package cron

import (
	"context"

	"github.com/isfk/xiaohui/config"
	cronV3 "github.com/robfig/cron/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewCron(lc fx.Lifecycle, log *zap.Logger, conf *config.Conf) *cronV3.Cron {
	c := cronV3.New(cronV3.WithSeconds())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting cron server")
			go c.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return c.Stop().Err()
		},
	})

	return c
}
