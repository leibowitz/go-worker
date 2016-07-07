package main

import (
	"flag"
	"fmt"
	"time"
	"runtime"
	"math/rand"
	"sync"

	"github.com/leibowitz/go-worker"
)

// Job represents the job to be run
type Job struct {
	Name string
	Wg   *sync.WaitGroup
	Rnd *rand.Rand
	JobQueue chan worker.Job
}

func (j Job) String() string {
	return j.Name
}

func (j Job) Process() error {
	fmt.Printf("Processing job %s\n", j.String())
	if j.Rnd.Intn(2) == 0 {
		return nil
	}
	return fmt.Errorf("Some kind of error")
}

func (j Job) Result(err error) {
	if err != nil {
		fmt.Printf("Processing job %s failed: %s\n", j.String(), err)
		j.JobQueue <- j
		return
	}
	fmt.Printf("Processed job %s successfully\n", j.String())
	j.Wg.Done()
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var jobs = flag.Int("j", 100, "Number of jobs")
	var workers = flag.Int("w", runtime.NumCPU(), "Number of workers")
	var queues = flag.Int("q", 2, "Number of queues")
	flag.Parse()

	dispatcher := worker.NewDispatcher(*workers)
	JobQueue := dispatcher.Run(*queues)

	wg := &sync.WaitGroup{}
	for i := 0; i < *jobs; i++ {
		wg.Add(1)
		fmt.Printf("Adding job %d to the queue\n", i)
		JobQueue <- Job{Name: fmt.Sprintf("%d", i), Wg: wg, Rnd: r, JobQueue: JobQueue}
	}

	wg.Wait()
	dispatcher.Stop()
}
