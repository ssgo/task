package task

import (
	"fmt"
	"github.com/ssgo/log"
	"github.com/ssgo/s"
	"github.com/ssgo/u"
	"time"
)

type Task struct {
	Group string
	Name  string
	Args  s.Map
}

type FetchedTask struct {
	Id   string
	Time int64
	Task Task
}

type ConfirmTask struct {
	Id string
	Ok bool
}

func (task *Task) PendingKey() string {
	return fmt.Sprint(task.Group, ".", task.Name)
}

func (task *Task) FailedKey() string {
	return fmt.Sprint(task.Group, ".", task.Name, "_failed")
}

func createTask(task Task) bool {
	if task.Name == "" {
		return false
	}

	logger.Info("add task", "task", task)
	return redisConn.LPUSH(task.PendingKey(), u.String(task.Args)) > 0
}

func fetchTask(task Task) *FetchedTask {
	if task.Name == "" {
		return nil
	}
	pendingKey := task.PendingKey()
	r := redisConn.RPOP(pendingKey)
	if r.String() == "" {
		return nil
	}
	task.Args = s.Map{}
	_ = r.To(&task.Args)

	fetchedTask := FetchedTask{Id: u.UniqueId(), Time: time.Now().Unix(), Task: task}
	logger.Info("fetch task", "fetchedTask", fetchedTask)

	// 将任务放入Doing队列
	redisConn.HSET("_task_doing", fetchedTask.Id, u.Json(fetchedTask))

	// 如果有重复的任务清除之
	n := redisConn.Do("LREM", pendingKey, 0, task.Args).Int()
	if n > 0 {
		logger.Warning(fmt.Sprint("removed ", n, " same tasks"), "task", task)
	}
	return &fetchedTask
}

func confirmTask(info ConfirmTask) bool {
	if info.Id == "" {
		return false
	}

	fetchedTask := FetchedTask{}
	_ = redisConn.HGET("_task_doing", info.Id).To(&fetchedTask)
	if fetchedTask.Id == "" {
		logger.Warning("task not exists", "fetchedTask", fetchedTask)
		return false
	}

	ok := redisConn.HDEL("_task_doing", info.Id) > 0
	if !ok {
		logger.Warning("task remove failed", "fetchedTask", fetchedTask)
		return false
	}

	startTime := time.Unix(fetchedTask.Time, 0)
	logger.Task(fetchedTask.Task.PendingKey(), fetchedTask.Task.Args, info.Ok, s.GetServerAddr(), startTime, log.MakeUesdTime(startTime, time.Now()), "ID: "+fetchedTask.Id)

	if !info.Ok {
		logger.Info("task failed", "fetchedTask", fetchedTask)

		// 将任务放入Failed队列
		return redisConn.LPUSH(fetchedTask.Task.FailedKey(), u.Json(fetchedTask)) > 0
	}

	logger.Info("task done", "fetchedTask", fetchedTask)
	return true
}
