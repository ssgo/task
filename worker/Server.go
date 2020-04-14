package worker

import (
	"github.com/ssgo/discover"
	"github.com/ssgo/log"
	"github.com/ssgo/s"
	"github.com/ssgo/u"
)

var logger = log.New(u.ShortUniqueId())
var caller *discover.Caller
var workers = make(map[string]func(task *FetchedTask) bool)

func Init() {
	loadConfig()
}

func Start() {
	s.NewTimerServer("TaskWorker", conf.CheckInterval.TimeDuration(), func(running *bool) {
		for _, taskName := range conf.Tasks {
			fetchedTask := FetchTask(taskName)
			if fetchedTask != nil {
				logger.Info("new work", "task", taskName, "fetchedTask", fetchedTask)
				f := workers[taskName]
				if f == nil {
					logger.Error("no worker for new work", "task", taskName, "fetchedTask", fetchedTask)
				} else {
					go func() {
						ok := f(fetchedTask)
						if ok {
							logger.Info("work succeed", "task", taskName, "fetchedTask", fetchedTask)
						} else {
							logger.Error("work failed", "task", taskName, "fetchedTask", fetchedTask)
						}
						ok2 := ConfirmTask(fetchedTask.Id, ok)
						if !ok2 {
							logger.Error("confirm failed", "task", taskName, "fetchedTask", fetchedTask)
						}
					}()
				}
			}
		}
	}, func() {
		caller = discover.NewCaller(nil, logger)
	}, nil)
}

func RegisterWorker(taskName string, f func(task *FetchedTask) bool) {
	workers[taskName] = f
}
