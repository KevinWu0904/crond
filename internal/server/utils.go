package server

import "encoding/json"

func unmarshalJob(data []byte) (*Job, error) {
	job := new(Job)
	err := json.Unmarshal(data, &job)
	return job, err
}

func marshalJob(job *Job) ([]byte, error) {
	return json.Marshal(job)
}
