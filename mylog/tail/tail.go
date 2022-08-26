package tail

import (
	"context"
	"fmt"
	"test/logagent/mylog/config"
	"test/logagent/mylog/kafka"
	"time"

	"github.com/hpcloud/tail"
)

type Tasks config.TailTask

var (
	tails      *tail.Tail
	tasks_map  map[string]*config.TailTask
	tasks_chan chan config.LogEntryConf
)

func run(T *config.TailTask) {
	for {
		select {
		case <-T.Ctx.Done():
			return
		case line := <-T.Instance.Lines:
			kafka.SendToChan(T.Topic, line.Text)
		}
	}
}

func Init(Tvalue config.LogEntryConf) error {
	tasks_map = make(map[string]*config.TailTask, 100)
	tasks_chan = make(chan config.LogEntryConf)
	for _, value := range Tvalue {
		base := config.LogConf{
			Path:  value.Path,
			Topic: value.Topic,
		}
		Task, err := NewTask(base)
		name := fmt.Sprintf("%s\\%s", value.Path, value.Topic)
		tasks_map[name] = Task
		if err != nil {
			fmt.Println("Init tail failed,err:", err)
			return err
		}
		go run(Task)
	}
	go Update_Task()
	return nil
}

func Update_Task() {
	for {
		select {
		case new_tasks := <-tasks_chan:
			for _, old_task := range tasks_map {
				name := fmt.Sprintf("%s\\%s", old_task.Path, old_task.Topic)
				tasks_map[name].CancelF()
			}
			for _, new_task := range new_tasks {
				name := fmt.Sprintf("%s\\%s", new_task.Path, new_task.Topic)
				Task, err := NewTask(*new_task)
				if err != nil {
					fmt.Println("init task err:", err)
					return
				}
				tasks_map[name] = Task
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func Get_chan() chan config.LogEntryConf {
	return tasks_chan
}

func NewTask(base config.LogConf) (tal *config.TailTask, err error) {
	cfg := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	ctx, cancel := context.WithCancel(context.Background())
	tails, err = tail.TailFile(base.Path, cfg)
	if err != nil {
		fmt.Println("tail file failed,err:", err)
	}
	tal = &config.TailTask{
		Path:     base.Path,
		Topic:    base.Topic,
		Instance: tails,
		Ctx:      ctx,
		CancelF:  cancel,
	}
	return
}
