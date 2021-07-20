package server

// Job represents crond Job entity in memory.
type Job struct {
	JobID          string       `json:"job_id"`
	JobKey         string       `json:"job_key"`
	JobDisplayName string       `json:"job_display_name"`
	CronExpression string       `json:"cron_expression"`
	ExecutorType   ExecutorType `json:"executor_type"`
}

// ExecutorType defines multiple executor types, different type will be running by different executors.
type ExecutorType int8

// Run implements cron.Job interface.
func (*Job) Run() {

}
