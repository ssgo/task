package task

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/ssgo/db"
	"github.com/ssgo/log"
	"github.com/ssgo/redis"
	"github.com/ssgo/s"
	"github.com/ssgo/u"
)

var logger = log.New(u.ShortUniqueId())
var redisConn *redis.Redis
var dbConn *db.DB

func Init() {
	s.Restful(1, "POST", "/{group}/{name}", createTask)
	s.Restful(2, "GET", "/{group}/{name}", fetchTask)
	s.Restful(2, "POST", "/confirm/{id}", confirmTask)

	loadConfig()
	redisConn = redis.GetRedis(conf.Redis, logger)
	dbConn = db.GetDB(conf.Db, logger)
}

func Start() {
	s.NewTimerServer("CronChecker", conf.CheckInterval.TimeDuration(), checkCron, func() {
		cronCron.Start()
	}, func() {
		cronCron.Stop()
	})
}
