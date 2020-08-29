///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package client

import (
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/queue"
	"time"
)

type job struct {
	finishChan chan error
	job        func() error
}

const BlockQueueSize = 100
const ProcessorSize = 10

var (
	workerBlockQueue *queue.NoneBlockQueue
)

func InitProcessors() {
	workerBlockQueue = queue.NewNoneBlockQueue(BlockQueueSize)
	for i := 0; i < ProcessorSize; i++ {
		go processor()
	}
	logger.Info("已启动", ProcessorSize, "个工作线程")
}

func Put(job *job) {
	for {
		success := workerBlockQueue.Put(job)
		if !success {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		break
	}
}

// 工作处理器
func processor() {
	for {
		q, s := workerBlockQueue.Fetch()
		if !s {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		if s && q != nil {
			j := q.(*job)
			gox.Try(func() {
				j.finishChan <- j.job()
			}, func(err interface{}) {
				j.finishChan <- err.(error)
			})
		}
	}
}
