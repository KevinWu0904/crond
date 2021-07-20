package server

// Job represents crond Job entity in memory.
type Job struct {
	JobID          string
	JobKey         string
	JobDisplayName string
	CronExpression string
	ExecutorType   ExecutorType
}

// ExecutorType defines multiple executor types, different type will be running by different executors.
type ExecutorType int8

// Run implements cron.Job interface.
func (*Job) Run() {

}
