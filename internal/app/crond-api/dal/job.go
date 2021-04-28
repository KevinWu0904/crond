package dal

type Job struct {
	JobID          int64
	JobName        string
	JobDescription string
	CronExpression string
}
