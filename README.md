# go-worker
Experiment with workers

Code from Marcio Castilho

http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/

# usage

```go
package main 

import (
  "runtime"
  "fmt"
  "time"
  
  "github.com/leibowitz/go-worker"
)

// The job type needs to implement the worker.Job interface
// type Job interface {
//        String() string
//        Process() error
//        Result(error)
// }
type Job struct {}

func (j Job) String() string {
  // return an identifier or something uniquely identifying this unit of work
  return ""
}

func (j Job) Process() error {
  // Do some computing here
  return nil
}

func (j Job) Result(err error) {
  // Handle success/failure in here (if needed)
  if err != nil {
    fmt.Print(err)
    return
  }
  
  fmt.Printf("Success %s\n", j)
}

func main() {
  dispatcher := worker.NewDispatcher(runtime.NumCPU()) // Create a worker per cpu
  JobQueue := dispatcher.Run(2) // Setup a dispatcher with a buffered channel of size 2
  
  for i := 0; i < 10; i++ {
    JobQueue <- Job{}
  }
  
  // Wait for a bit
  time.Sleep(1*time.Second) // Better to use sync.WaitGroup to wait for all processes to finish
  
  // Kill all workers
  dispatcher.Stop()
}
```
