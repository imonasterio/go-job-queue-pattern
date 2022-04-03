package worker

import (
	"fmt"

	"github.com/imonasterio/go-job-queue-pattern/cmd/job"
)

type Worker struct {
	ID             int
	jobs           chan *job.Job
	dispatchStatus chan *DispatchStatus
	Quit           chan bool
}

func CreateNewWorker(id int, workerQueue chan *Worker, jobQueue chan *job.Job, dStatus chan *DispatchStatus) *Worker {
	w := &Worker{
		ID:             id,
		jobs:           jobQueue,
		dispatchStatus: dStatus,
	}

	go func() { workerQueue <- w }()
	return w
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case job := <-w.jobs:
				fmt.Printf("Worker[%d] executing job[%d].\n", w.ID, job.ID)
				job.F()
				w.dispatchStatus <- &DispatchStatus{Type: "worker", ID: w.ID, Status: "quit"}
				w.Quit <- true
			case <-w.Quit:
				return
			}
		}
	}()
}
