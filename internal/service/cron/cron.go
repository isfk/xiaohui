package cron

import (
	cronV3 "github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/isfk/xiaohui/config"
	"github.com/isfk/xiaohui/internal/service/cron/jobs"
)

type CronJobs struct{}

func NewCronJobs(c *cronV3.Cron, log *zap.Logger, conf *config.Conf, audio *jobs.AudioJob) *CronJobs {
	list := map[string]Job{
		audio.String(): audio,
	}

	log.Info("cron", zap.Any("config", conf.Config.Cron.Jobs))

	for _, job := range conf.Config.Cron.Jobs {
		if job.Enabled {
			id, err := c.AddJob(job.Spec, list[job.Name])
			if err != nil {
				log.With(zap.Error(err)).Error("Job add fail", zap.String("name", job.Name), zap.Any("id", id))
			} else {
				log.Info("Job add success, id: %v", zap.String("name", job.Name), zap.Any("id", id))
			}
		}
	}
	return &CronJobs{}
}
