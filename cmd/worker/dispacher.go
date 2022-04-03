package worker

import (
	"fmt"
	"sync"

	"github.com/imonasterio/go-job-queue-pattern/cmd/job"
)

type DispatchStatus struct {
	Type   string
	ID     int
	Status string
}

type Dispatcher struct {
	jobCounter     int                  // internal counter for number of jobs
	jobQueue       chan *job.Job        // channel of jobs submitted by main()
	dispatchStatus chan *DispatchStatus // channel for job/worker status reports
	workQueue      chan *job.Job        // channel of work dispatched
	workerQueue    chan *Worker         // channel of workers
	sync.RWMutex
}

func CreateNewDispatcher() *Dispatcher {
	d := &Dispatcher{
		jobCounter:     0,
		jobQueue:       make(chan *job.Job),
		dispatchStatus: make(chan *DispatchStatus),
		workQueue:      make(chan *job.Job),
		workerQueue:    make(chan *Worker),
		RWMutex:        sync.RWMutex{},
	}
	return d
}

func (d *Dispatcher) Start(numWorkers int) {
	// Create numWorkers:
	for i := 0; i < numWorkers; i++ {
		worker := CreateNewWorker(i, d.workerQueue, d.workQueue, d.dispatchStatus)
		worker.Start()
	}

	// wait for work to be added then pass it off.
	go func() {
		for {
			select {
			case job := <-d.jobQueue:
				fmt.Printf("Got a job in the queue to dispatch: %d\n", job.ID)
				// Sending it off;
				d.workQueue <- job
			case ds := <-d.dispatchStatus:
				fmt.Printf("Got a dispatch status:\n\tType[%s] - ID[%d] - Status[%s]\n", ds.Type, ds.ID, ds.Status)
				d.Lock()
				if ds.Type == "worker" {
					if ds.Status == "quit" {
						d.jobCounter--
					}
				}
				d.Unlock()
			}
		}
	}()
}

func (d *Dispatcher) AddJob(je job.JobExecutable) {
	j := &job.Job{ID: d.jobCounter, F: je}
	go func() { d.jobQueue <- j }()
	d.jobCounter++
	fmt.Printf("Quantity jobs is now: %d\n", d.jobCounter)
}

func (d *Dispatcher) Finished() bool {
	d.RLock()
	defer d.RUnlock()
	if d.jobCounter < 1 {
		return true
	} else {
		return false
	}
}
