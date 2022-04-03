package cli

import (
	"fmt"
	"time"

	"github.com/imonasterio/go-job-queue-pattern/cmd/worker"
)

func Start() {
	job1 := func() error {
		fmt.Printf("job1: doing some task.\n")
		time.Sleep(2 * time.Second)
		fmt.Printf("[job1] Work done.\n")
		return nil
	}

	job2 := func() error {
		fmt.Printf("job2: doing some task..\n")
		time.Sleep(4 * time.Second)
		fmt.Printf("[job2] Work done.\n")
		return nil
	}

	d := worker.CreateNewDispatcher()
	d.AddJob(job1)
	d.AddJob(job2)
	d.Start(2)

	for {
		if d.Finished() {
			fmt.Printf("All jobs finished.\n")
			break
		}
	}
}
