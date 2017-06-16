package models

import "time"

const (
	BuildStatusQueue   = "queued"
	BuildStatusStart   = "started"
	BuildStatusCancel  = "canceled"
	BuildStatusFail    = "failed"
	BuildStatusSuccess = "succeeded"
)

type BuildStatus struct {
	QueueTimestamp    int64
	StartTimestamp    int64
	CompleteTimestamp int64
	Status            string
}

func NewBuildStatus() *BuildStatus {
	return &BuildStatus{}
}

func (b *BuildStatus) SetQueued() {
	b.Status = BuildStatusQueue
	b.QueueTimestamp = time.Now().Unix()
}

func (b *BuildStatus) SetStart() {
	b.Status = BuildStatusStart
	b.StartTimestamp = time.Now().Unix()
}

func (b *BuildStatus) SetSuccess() {
	b.Status = BuildStatusSuccess
	b.CompleteTimestamp = time.Now().Unix()
}
