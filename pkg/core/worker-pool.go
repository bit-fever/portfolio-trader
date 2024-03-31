//=============================================================================
/*
Copyright © 2024 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
//=============================================================================

package core

import (
	"log/slog"
	"time"
)

//=============================================================================

type WorkerPool struct {
	numWorkers int
	queueSize  int
	taskQueue  chan func()
}

//=============================================================================
//===
//=== Methods
//===
//=============================================================================

func (p *WorkerPool) Init(numWorkers, queueSize int) {
	p.numWorkers = numWorkers
	p.queueSize  = queueSize
	p.taskQueue  = make(chan func(), queueSize)

	for i := 0; i<numWorkers; i++ {
		go p.worker()
	}

	slog.Info("Init: Worker pool created", "workers", numWorkers)
}

//=============================================================================

func (p *WorkerPool) ShutDown() bool {

	if p.taskQueue != nil {
		close(p.taskQueue)

		//--- Wait for task completion
		for ; len(p.taskQueue) > 0 ; {
			time.Sleep(time.Millisecond * 500)
		}

		p.taskQueue = nil
	}

	slog.Info("ShutDown: Worker pool stopped")

	return true
}

//=============================================================================

func (p *WorkerPool) Submit(task func()) {
	p.taskQueue <- task
}

//=============================================================================
//===
//=== Worker
//===
//=============================================================================

func (p *WorkerPool) worker() {
	for {
		select {
			case task, ok := <- p.taskQueue:
				if ok {
					//--- Run task
					task()
				} else {
					//--- Channel closed. Exit from goroutine
					return
				}

			default:
				time.Sleep(time.Millisecond * 50)
		}
	}
}

//=============================================================================
