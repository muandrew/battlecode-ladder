package build

import (
	"github.com/jeffail/tunny"
	"github.com/muandrew/battlecode-ladder/utils"
	"strconv"
)

type worker struct {
	id int
}

func CreateWorkers(numWorkers int) []tunny.TunnyWorker {
	workers := make([]tunny.TunnyWorker, numWorkers)
	for i := range workers {
		workers[i] = &worker{
			id: i,
		}
		utils.RunShell("sh", []string{"scripts/setup-worker-match-workspace.sh", strconv.Itoa(i)})
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
