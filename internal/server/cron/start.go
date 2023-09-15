package cron

import (
	cronV3 "github.com/robfig/cron/v3"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/isfk/xiaohui/config"
	"github.com/isfk/xiaohui/internal/pkg/logger"
	"github.com/isfk/xiaohui/internal/service/cron"
	"github.com/isfk/xiaohui/internal/service/cron/jobs"
)

func Start() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			config.NewConfig,
			NewCron,
			jobs.NewAudioJob,
			cron.NewCronJobs,
			func(conf *config.Conf) *zap.Logger {
				return logger.NewLogger("cron", conf.Config.Cron.Log)
			},
		),
		fx.Invoke(func(h *cronV3.Cron) {}, func(h *cron.CronJobs) {}),
	).Run()
}
