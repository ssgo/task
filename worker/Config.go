package worker

import (
	"fmt"
	"github.com/ssgo/config"
	"time"
)

var conf = struct {
	ServerApp     string
	Tasks         []string
	CheckInterval config.Duration
}{}

func loadConfig() {
	config.LoadConfig("worker", &conf)
	fmt.Println(conf)
	if conf.CheckInterval == 0 {
		conf.CheckInterval = config.Duration(10 * time.Second)
	}
	if conf.ServerApp == "" {
		conf.ServerApp = "task"
	}
}
