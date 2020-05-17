package models

import "time"

const (
	//BuildStatusQueue queued
	BuildStatusQueue = "queued"
	//BuildStatusStart started
	BuildStatusStart = "started"
	//BuildStatusCancel canceled
	BuildStatusCancel = "canceled"
	//BuildStatusFail failed
	BuildStatusFail = "failed"
	//BuildStatusSuccess success
	BuildStatusSuccess = "succeeded"
)

//BuildStatus represents an event in the build process
type BuildStatus struct {
	QueueTimestamp    int64
	StartTimestamp    int64
	CompleteTimestamp int64
	Status            string
}

//NewBuildStatus creates a new instance of BuildStatus
func NewBuildStatus() *BuildStatus {
	return &BuildStatus{}
}

//SetQueued sets the status to Queued, also sets the time.
func (b *BuildStatus) SetQueued() {
	b.Status = BuildStatusQueue
	b.QueueTimestamp = time.Now().Unix()
}

//SetStart sets the status to Start, also sets the time.
func (b *BuildStatus) SetStart() {
	b.Status = BuildStatusStart
	b.StartTimestamp = time.Now().Unix()
}

//SetSuccess sets the status to Success, also sets the time.
func (b *BuildStatus) SetSuccess() {
	b.Status = BuildStatusSuccess
	b.CompleteTimestamp = time.Now().Unix()
}

//SetFailure sets the status to Failure, also sets the time.
func (b *BuildStatus) SetFailure() {
	b.Status = BuildStatusFail
	b.CompleteTimestamp = time.Now().Unix()
}
