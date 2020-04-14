# a simple task server

### configure a db source for set crontab

crontab will auto create task to queue

if task never be consumed will be over 

```mysql
CREATE TABLE `Crontab` (
  `id` int NOT NULL AUTO_INCREMENT,
  `group` varchar(20) NOT NULL,
  `name` varchar(20) NOT NULL,
  `spec` varchar(100) NOT NULL,
  `args` varchar(10240) NOT NULL,
  `active` enum('true','false') NOT NULL DEFAULT 'true',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;
```

crontab setting reference by https://github.com/robfig/cron

thanks for Rob Figueiredo

### how to create a task by api?

looking for https://github.com/ssgo/task/worker


### how to fetch a task by api?

looking for https://github.com/ssgo/task/worker
