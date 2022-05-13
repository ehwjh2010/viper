package routine

import (
	"errors"
	"runtime/debug"
	"sync"
	"sync/atomic"

	"github.com/panjf2000/ants/v2"

	"github.com/ehwjh2010/viper/client/settings"
)

type Task struct {
	rawConfig         settings.Routine
	p                 *ants.Pool
	rawGoroutineCount int64
}

var (
	backgroundTask *Task
	mtx            sync.Mutex

	NoEnableRoutinePool = errors.New("not enable routine pool, please add enableRtePool field to settings")
)

func newTask(rawConfig settings.Routine, p *ants.Pool) *Task {
	return &Task{rawConfig: rawConfig, p: p}
}

// Close 关闭协程池
func (task *Task) Close() {
	if task == nil || task.p == nil || task.p.IsClosed() {
		return
	}

	task.p.Release()
}

// Reboot 重启关闭的协程池
func (task *Task) Reboot() {
	if task == nil || task.p == nil || task.p.IsClosed() {
		return
	}

	task.p.Reboot()
}

// incr 原生协程数量+1
func (task *Task) incr() {
	atomic.AddInt64(&task.rawGoroutineCount, 1)
}

// decr 原生协程数量-1
func (task *Task) decr() {
	atomic.AddInt64(&task.rawGoroutineCount, -1)
}

// wrapper 包装任务函数
func (task *Task) wrapperForRaw(f TaskFunc) TaskFunc {
	return func() {
		defer task.decr()
		defer func() {
			if e := recover(); e != nil {
				task.rawConfig.Logger.Printf("raw goroutine execute task occur panic, err: %s,  stack ==> %s", e, string(debug.Stack()))
			}
		}()

		task.incr()
		f()
	}
}

// RawGoroutineCount 获取当前执行任务的原生协程数
func (task *Task) RawGoroutineCount() int64 {
	return task.rawGoroutineCount
}

// AsyncDO 添加任务, 如果有设置
func (task *Task) AsyncDO(taskFunc TaskFunc) error {
	err := task.p.Submit(taskFunc)

	if err == nil {
		return nil
	}

	if errors.Is(err, ants.ErrPoolOverload) {
		if task.rawConfig.UseRawWhenBusy {
			go task.wrapperForRaw(taskFunc)()
			return nil
		}

		return err
	}

	return err

}

// CountInfo 获取协程池个数信息
func (task Task) CountInfo() map[string]int {
	result := make(map[string]int, 3)

	result["free"] = task.p.Free()
	result["running"] = task.p.Running()
	result["cap"] = task.p.Cap()

	return result
}

// SetUpDefaultTask 初始化后台任务
func SetUpDefaultTask(conf settings.Routine) func() error {
	return func() error {
		if backgroundTask != nil {
			return nil
		}

		mtx.Lock()
		defer mtx.Unlock()
		if backgroundTask != nil {
			return nil
		}

		if task, err := SetUp(conf); err != nil {
			return err
		} else {
			backgroundTask = task
		}
		return nil

	}
}

func CloseDefaultTask() error {
	if backgroundTask == nil {
		return nil
	}

	backgroundTask.Close()
	return nil
}

// AddTask 添加任务
func AddTask(taskFunc TaskFunc) error {
	if backgroundTask == nil {
		return NoEnableRoutinePool
	}

	return backgroundTask.AsyncDO(taskFunc)
}

// CountInfo 获取协程池个数信息
func CountInfo() (map[string]int, error) {
	if backgroundTask == nil {
		return nil, NoEnableRoutinePool
	}

	countInfo := backgroundTask.CountInfo()
	return countInfo, nil
}
