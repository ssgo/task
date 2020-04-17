package worker

import (
	"github.com/ssgo/s"
	"github.com/ssgo/u"
	"strings"
)

type Task struct {
	Group string
	Name  string
	Args  map[string]interface{}
}

type FetchedTask struct {
	Id   string
	Time int64
	Task Task
}

func (task *Task) ArgsTo(to interface{}) {
	u.Convert(task.Args, to)
}

func FetchTask(taskName string) *FetchedTask {
	if taskName == "" {
		return nil
	}
	fetchedTask := new(FetchedTask)
	r := caller.Get(conf.ServerApp, "/"+strings.Replace(taskName, ".", "/", 1))
	if r.Error != nil {
		logger.Error(r.Error.Error())
	}
	_ = r.To(fetchedTask)
	if fetchedTask.Id == "" {
		return nil
	}
	return fetchedTask
}

func ConfirmTask(id string, ok bool) bool {
	if id == "" {
		return false
	}
	r := caller.Post(conf.ServerApp, "/confirm/"+id, s.Map{"ok": ok})
	if r.Error != nil {
		logger.Error(r.Error.Error())
	}
	return r.String() == "true"
}
