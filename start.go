package main

import (
	"github.com/ssgo/s"
	"github.com/ssgo/task/task"
)

func main() {
	task.Init()
	task.Start()
	s.Start()
}
