package server

import (
	"context"
	"sync"

	"github.com/KevinWu0904/crond/internal/common/constant"
	"github.com/KevinWu0904/crond/pkg/logs"

	"github.com/robfig/cron/v3"
)

// CronDispatcher represents crond unified job dispatcher which is expected to be running only in raft leader node.
type CronDispatcher struct {
	Cron       *cron.Cron
	JobEntries sync.Map

	sync.Mutex

	started bool
}

// NewCronDispatcher creates CronDispatcher.
func NewCronDispatcher() *CronDispatcher {
	return &CronDispatcher{
		Cron: cron.New(cron.WithSeconds()),
	}
}

// Start will load initial jobs from persistent storage and start CronDispatcher.
func (cd *CronDispatcher) Start(ctx context.Context, initJobs []*Job) error {
	cd.Lock()
	defer cd.Unlock()

	if cd.started {
		return nil
	}

	for _, job := range initJobs {
		if err := cd.AddJob(ctx, job); err != nil {
			return err
		}
	}

	cd.Cron.Start()
	cd.started = true

	return nil
}

// Stop will try to stop CronDispatcher gracefully and will shutdown CronDispatcher forcibly after timeout.
func (cd *CronDispatcher) Stop(ctx context.Context) error {
	cd.Lock()
	defer cd.Unlock()

	if !cd.started {
		return nil
	}

	stop := cd.Cron.Stop()
	cd.started = false

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-stop.Done():
	}

	return nil
}

// AddJob adds a new Job into existing CronDispatcher.
func (cd *CronDispatcher) AddJob(ctx context.Context, job *Job) error {
	ctx = logs.CtxAddKVs(ctx, constant.LogJobKey, job.JobKey)

	if _, ok := cd.JobEntries.Load(job.JobKey); ok {
		cd.DeleteJob(ctx, job)
	}

	entryID, err := cd.Cron.AddJob(job.CronExpression, job)
	if err != nil {
		logs.CtxError(ctx, "AddJob failed: err=%v", err)
		return err
	}

	cd.JobEntries.Store(job.JobKey, entryID)
	logs.CtxInfo(ctx, "AddJob successfully: entryID=%d", entryID)
	return nil
}

// DeleteJob deletes an existing job from CronDispatcher.
func (cd *CronDispatcher) DeleteJob(ctx context.Context, job *Job) {
	ctx = logs.CtxAddKVs(ctx, constant.LogJobKey, job.JobKey)

	if entryID, ok := cd.JobEntries.Load(job.JobKey); ok {
		cd.Cron.Remove(entryID.(cron.EntryID))
		cd.JobEntries.Delete(job.JobKey)

		logs.CtxInfo(ctx, "DeleteJob successfully: entryID=%d", entryID)
	}
}
