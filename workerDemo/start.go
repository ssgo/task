package main

import (
	"github.com/ssgo/log"
	"github.com/ssgo/s"
	"github.com/ssgo/task/worker"
	"github.com/ssgo/u"
)

func main() {
	worker.RegisterWorker("message.sms", sendSms)
	worker.Init()
	worker.Start()
	s.Start()
}

func sendSms(task *worker.FetchedTask) bool {
	if u.GlobalRand1.Int()%3 == 1 {
		log.DefaultLogger.Error("send failed", "task", task)
		return false
	} else {
		log.DefaultLogger.Info("send succeed", "task", task)
		return true
	}
}

/*
creat test data for make message.sms job every 10s:

CREATE TABLE `Crontab` (
  `id` int NOT NULL AUTO_INCREMENT,
  `group` varchar(20) NOT NULL,
  `name` varchar(20) NOT NULL,
  `spec` varchar(100) NOT NULL,
  `args` varchar(10240) NOT NULL,
  `active` enum('true','false') NOT NULL DEFAULT 'true',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;

INSERT INTO `Crontab` (`group`,`name`,`spec`,`args`,`active`) VALUES ('message','sms','@every 10s','{"message":"test","phone":"13838384384"}',true);
 */