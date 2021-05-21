package dal

import (
	"context"

	"github.com/KevinWu0904/crond/pkg/logs"
	"github.com/KevinWu0904/crond/pkg/mysql"
)

type Job struct {
	JobID          string
	JobName        string
	JobDescription string
	CronExpression string
}

func (Job) TableName() string {
	return "crond_job"
}

func CreateJob(ctx context.Context, job *Job) error {
	if err := mysql.Client.WithContext(ctx).Create(job).Error; err != nil {
		logs.CtxError(ctx, "CreateJob failed: err=%v", err)
		return err
	}

	return nil
}
