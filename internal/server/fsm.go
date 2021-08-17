package server

import (
	"io"

	"github.com/hashicorp/raft"
)

// JobFSM can be used as a Raft.FSM that can be implemented by
// clients to make use of the replicated log.
type JobFSM struct {
	crondStore CrondStore
}

// NewJobFSM constructs a JobFSM with crondStore as job storage
func NewJobFSM(crondStore CrondStore) JobFSM {
	return JobFSM{
		crondStore: crondStore,
	}
}

// Apply log is invoked once a log entry is committed.
// It returns response of Raft.Apply
func (f JobFSM) Apply(*raft.Log) interface{} {
	panic("implement me")
}

// Snapshot returns an Snapshot which can be used to save a point-in-time snapshot of the FSM
func (f JobFSM) Snapshot() (raft.FSMSnapshot, error) {
	panic("implement me")
}

// Restore is used to restore an FSM from a snapshot
func (f JobFSM) Restore(io.ReadCloser) error {
	panic("implement me")
}

// Snapshot is returned by an FSM in response to a Snapshot
type Snapshot struct {
}

// Persist should dump all necessary state to the WriteCloser 'sink',
// and call sink.Close() when finished or call sink.Cancel() on error.
func (s Snapshot) Persist(sink raft.SnapshotSink) error {
	panic("implement me")
}

// Release is invoked when we are finished with the snapshot.
func (s Snapshot) Release() {
	panic("implement me")
}
