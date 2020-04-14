package task

import (
	"github.com/ssgo/config"
	"time"
)

var conf = struct {
	Cron          map[string]string
	CheckInterval config.Duration
	Redis         string
	Db            string
}{}

func loadConfig() {
	config.LoadConfig("task", &conf)
	if conf.CheckInterval == 0 {
		conf.CheckInterval = config.Duration(10 * time.Second)
	}
	if conf.Redis == "" {
		conf.Redis = "redis://127.0.0.1:6379/14"
	}
	if conf.Db == "" {
		conf.Db = "mysql://root:@127.0.0.1/task"
	}
}
