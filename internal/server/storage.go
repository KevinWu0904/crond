package server

import (
	"errors"
	"github.com/boltdb/bolt"
	"sync"
)

const (
	// Permissions to use on the db file. This is only used if the
	// database file does not exist and needs to be created.
	dbFileMode = 0600
)

var (
	// Bucket names we perform transactions in
	dbJobs = []byte("jobs")

	// An error indicating a given key does not exist
	ErrKeyNotFound = errors.New("not found")
)

// CrondStorage is used to provide an interface for storing
// and retrieving jobs in a durable fashion.
type CrondStorage interface {
	// SetJob creates or updates a job.
	SetJob(job *Job) error

	// DeleteJob deletes a job by key.
	DeleteJob(key string) error

	// GetJobByKey returns the value for key, or an empty job if key was not found.
	GetJobByKey(key string) (*Job, error)

	// GetJobs returns all jobs.
	GetJobs() ([]*Job, error)
}

// CrondBoltStore provides access to BoltDB for Crond to store and retrieve
// job entries. It also provides key/value storage, and can be used as a CrondStorage .
type CrondBoltStore struct {
	conn *bolt.DB
	path string
}

// NewCrondBoltStore takes a file path and returns a connected bolt store.
func NewCrondBoltStore(path string) (*CrondBoltStore, error) {
	handle, err := bolt.Open(path, dbFileMode, nil)
	if err != nil {
		return nil, err
	}
	return &CrondBoltStore{
		conn: handle,
		path: path,
	}, nil
}

// SetJob creates or updates a job on bolt db.
func (s *CrondBoltStore) SetJob(job *Job) error {
	tx, err := s.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	val, err := marshalJob(job)
	if err != nil {
		return err
	}
	bucket := tx.Bucket(dbJobs)
	if err := bucket.Put([]byte(job.JobKey), val); err != nil {
		return err
	}
	return tx.Commit()
}

// DeleteJob deletes a job by key on bolt db.
func (s *CrondBoltStore) DeleteJob(key string) error {
	tx, err := s.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(dbJobs)
	if err := bucket.Delete([]byte(key)); err != nil {
		return err
	}
	return tx.Commit()
}

// GetJobByKey returns the value for key, or an empty job if key was not found on bold db.
func (s *CrondBoltStore) GetJobByKey(key string) (*Job, error) {
	tx, err := s.conn.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(dbJobs)
	val := bucket.Get([]byte(key))
	if val == nil {
		return nil, ErrKeyNotFound
	}
	return unmarshalJob(val)
}

// GetJobs returns all jobs on bolt db.
func (s *CrondBoltStore) GetJobs() ([]*Job, error) {
	tx, err := s.conn.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(dbJobs)
	c := bucket.Cursor()
	jobs := make([]*Job, 0)
	for k, v := c.First(); k != nil; k, v = c.Next() {
		job, err := unmarshalJob(v)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// CrondInmemStore implements the CrondStorage interface.
// It should NOT EVER be used for production. It is used only for
// unit tests. Use the MDBStore implementation instead.
type CrondInmemStore struct {
	mu     sync.RWMutex
	jobMap map[string]*Job
}

// NewCrondInmemStore returns a mem job storage
func NewCrondInmemStore() *CrondInmemStore {
	return &CrondInmemStore{
		jobMap: make(map[string]*Job),
	}
}

// SetJob creates or updates a job in mem.
func (s *CrondInmemStore) SetJob(job *Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobMap[job.JobKey] = job
	return nil
}

// DeleteJob deletes a job by key in mem.
func (s *CrondInmemStore) DeleteJob(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobMap, key)
	return nil
}

// GetJobByKey returns the value for key, or an empty job if key was not found in mem.
func (s *CrondInmemStore) GetJobByKey(key string) (*Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if job, ok := s.jobMap[key]; ok {
		return job, nil
	}
	return nil, ErrKeyNotFound
}

// GetJobs returns all jobs in mem.
func (s *CrondInmemStore) GetJobs() ([]*Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	jobs := make([]*Job, len(s.jobMap))
	for _, job := range s.jobMap {
		jobs = append(jobs, job)
	}
	return jobs, nil
}
