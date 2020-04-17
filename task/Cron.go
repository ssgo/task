package task

import (
	"encoding/json"
	"github.com/robfig/cron/v3"
	"github.com/ssgo/s"
	"github.com/ssgo/u"
	"time"
)

type CronJob struct {
	entryID cron.EntryID
	hash    string
	Group   string
	Name    string
	Spec    string
	Args    string
	Active  bool
	//args    map[string]interface{}
}

func (job *CronJob) Run() {
	args := make(map[string]interface{})
	err := json.Unmarshal([]byte(job.Args), &args)
	if err != nil {
		logger.Error(err.Error(), "job", job)
	}
	createTask(Task{Group: job.Group, Name: job.Name, Args: args})
}

var cronJobs = make(map[int]*CronJob)
var cronCron = cron.New()

func checkCron(running *bool) {
	// 获取更新数据的事务锁
	expires := int(conf.CheckInterval.TimeDuration()/time.Second) * 3
	redisConn.Do("SET", "_cron_lock", s.GetServerAddr(), "NX", "EX", expires).Bool()
	lockedNode := redisConn.GET("_cron_lock").String()
	//logger.Info("get cron lock", "lockedNode", lockedNode, "myNode", s.GetServerAddr(), "succeed", lockedNode == s.GetServerAddr(), "cmd", fmt.Sprintln("SET", "_cron_syncing", s.GetServerAddr(), "NX", "PX", expires))
	if lockedNode == s.GetServerAddr() {
		redisConn.EXPIRE("_cron_lock", expires)
		// 设置标识位
		for _, job := range cronJobs {
			job.Active = false
		}

		// 载入 crontab
		dbConn.Query("SELECT `id`, `group`, `name`, `spec`, `args`, `active` FROM `Crontab` WHERE `active`='true'").ToKV(&cronJobs)
		if !*running {
			return
		}

		// 更新数据
		for jobId, job := range cronJobs {
			if !*running {
				return
			}

			//job.args = make(map[string]interface{})
			//err := json.Unmarshal([]byte(job.Args), &job.args)
			//if err != nil {
			//	logger.Error(err.Error(), "job", job)
			//}

			if job.Active {
				hash := u.Sha1String(u.Json(job))
				if hash != job.hash {
					// 需要更新
					var err error
					var newJobId cron.EntryID
					newJobId, err = cronCron.AddJob(job.Spec, job)
					if err != nil {
						logger.Error(err.Error(), "job", job)
						continue
					}

					// 清除旧Job
					if job.entryID != 0 {
						cronCron.Remove(job.entryID)
						logger.Info("update job", "job", job)
					} else {
						logger.Info("add job", "job", job)
					}

					job.hash = hash
					job.entryID = newJobId
				}
			} else {
				// 删除不存在的
				if job.entryID != 0 {
					cronCron.Remove(job.entryID)
				}
				delete(cronJobs, jobId)
				logger.Info("remove job", "job", job)
			}
		}
	} else {
		// 清除非cron工作节点的定时发生器
		if len(cronJobs) > 0 {
			for _, job := range cronJobs {
				if job.entryID != 0 {
					cronCron.Remove(job.entryID)
				}
				job.entryID = 0
				logger.Info("remove job", "job", job)
			}
			cronJobs = make(map[int]*CronJob)
		}
	}
}
