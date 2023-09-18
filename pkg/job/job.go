package job

import "time"

type Job interface {
	Run()
	SleepTime() time.Duration
}

func RegisterJob(j Job) {
	go func(j Job) {
		for {
			j.Run()
			time.Sleep(j.SleepTime())
		}
	}(j)
}
