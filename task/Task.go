package task

import (
	"fmt"
	"github.com/ssgo/u"
	"time"
)

type Task struct {
	Group string
	Name  string
	Args  string
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
	if task.Args == "" {
		task.Args = "{}"
	}
	logger.Info("add task", "task", task)
	return redisConn.LPUSH(task.PendingKey(), task.Args) > 0
}

func fetchTask(task Task) *FetchedTask {
	if task.Name == "" {
		return nil
	}
	pendingKey := task.PendingKey()
	task.Args = redisConn.RPOP(pendingKey).String()
	if task.Args == "" {
		return nil
	}

	fetchedTask := FetchedTask{Id: u.UniqueId(), Time: time.Now().Unix(), Task: task}
	logger.Info("fetch task", "task", fetchedTask)

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
	redisConn.HGET("_task_doing", info.Id).To(&fetchedTask)
	if fetchedTask.Id == "" {
		logger.Warning("task not exists", "id", info.Id, "ok", info.Ok)
		return false
	}

	ok := redisConn.HDEL("_task_doing", info.Id) > 0
	if !ok {
		logger.Warning("task remove failed", "id", info.Id, "ok", info.Ok, "time", fetchedTask.Time, "task", fetchedTask.Task)
		return false
	}

	if !info.Ok {
		logger.Info("task failed", "id", info.Id, "ok", info.Ok, "time", fetchedTask.Time, "task", fetchedTask.Task)

		// 将任务放入Failed队列
		return redisConn.LPUSH(fetchedTask.Task.FailedKey(), u.Json(fetchedTask)) > 0
	}

	logger.Info("task done", "id", info.Id, "ok", info.Ok, "time", fetchedTask.Time, "task", fetchedTask.Task)
	return true
}
