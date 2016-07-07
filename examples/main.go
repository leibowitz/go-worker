package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"

	"github.com/leibowitz/go-worker"
)

// Job represents the job to be run
type Job struct {
	Name string
	Wg   *sync.WaitGroup
}

func (j Job) String() string {
	return j.Name
}

func (j Job) Process() error {
	fmt.Printf("Processing job %s\n", j.String())
	defer j.Wg.Done()
	return nil
}

func main() {
	var jobs = flag.Int("j", 100, "Number of jobs")
	var workers = flag.Int("w", runtime.NumCPU(), "Number of workers")
	var queues = flag.Int("q", 2, "Number of queues")
	flag.Parse()

	JobQueue := make(chan worker.Job, *queues)
	dispatcher := worker.NewDispatcher(*workers)
	dispatcher.Run(JobQueue)

	wg := &sync.WaitGroup{}
	for i := 0; i < *jobs; i++ {
		wg.Add(1)
		fmt.Printf("Adding job %d to the queue\n", i)
		JobQueue <- Job{Name: fmt.Sprintf("%d", i), Wg: wg}
	}

	wg.Wait()
	dispatcher.Stop()
}
