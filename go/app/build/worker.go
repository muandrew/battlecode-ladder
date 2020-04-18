package build

import (
	"github.com/jeffail/tunny"
)

type worker struct {
	id int
}

func CreateWorkers(workerDir string, numWorkers int) []tunny.TunnyWorker {
	workers := make([]tunny.TunnyWorker, numWorkers)
	for i := range workers {
		workers[i] = &worker{
			id: i,
		}
	}
	return workers
}

func (w *worker) TunnyJob(data interface{}) interface{} {
	method := data.(func(int))
	method(w.id)
	return nil
}

func (w *worker) TunnyReady() bool {
	return true
}
