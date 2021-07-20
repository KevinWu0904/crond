package server

import (
	"net"
	"path"
	"time"

	"github.com/KevinWu0904/crond/pkg/logs"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

const (
	raftFileSnapshotStoreRetain = 3
	raftMaxLogCacheSize         = 500
	raftNetworkTransportMaxPool = 3
	raftNetworkTransportTimeout = time.Second * 30
)

// RaftStreamLayer implements raft low-level network transport.
type RaftStreamLayer struct {
	net.Listener
}

// NewRaftStreamLayer creates RaftStreamLayer.
func NewRaftStreamLayer(listener net.Listener) *RaftStreamLayer {
	return &RaftStreamLayer{
		listener,
	}
}

// Dial connects to the address on the named network.
func (t *RaftStreamLayer) Dial(address raft.ServerAddress, timeout time.Duration) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: timeout}
	return dialer.Dial("tcp", string(address))
}

// RaftLayer represents crond raft consensus layer.
type RaftLayer struct {
	bootstrap bool
	underlay  *raft.Raft

	rc            *raft.Config
	snapshotStore raft.SnapshotStore
	stableStore   raft.StableStore
	logStore      raft.LogStore
	transport     raft.Transport
}

// NewRaftLayer creates crond RaftLayer.
func NewRaftLayer(c *Config, listener net.Listener) *RaftLayer {
	rc := raft.DefaultConfig()
	rc.LogOutput = logs.GetRaftWriter()
	rc.LocalID = raft.ServerID(c.RaftNode)

	var err error

	var snapshotStore raft.SnapshotStore
	var stableStore raft.StableStore
	var logStore raft.LogStore
	var jobStore Storage

	if !c.RaftProduction {
		snapshotStore = raft.NewInmemSnapshotStore()

		memStore := raft.NewInmemStore()
		stableStore = memStore
		logStore = memStore
		jobStore = NewJobInmemStore()
	} else {
		snapshotStore, err = raft.NewFileSnapshotStore(path.Join(c.RaftDataDir, "raft"), raftFileSnapshotStoreRetain, logs.GetRaftWriter())
		if err != nil {
			logs.Fatal("NewRaftLayer failed to create snapshotStore: err=%v", err)
		}

		boltStore, err := raftboltdb.NewBoltStore(path.Join(c.RaftDataDir, "raft", "raft.db"))
		if err != nil {
			logs.Fatal("NewRaftLayer failed to create stable store: err=%v", err)
		}

		stableStore = boltStore

		logStore, err = raft.NewLogCache(raftMaxLogCacheSize, boltStore)
		if err != nil {
			logs.Fatal("NewRaftLayer failed to create log store: err=%v", err)
		}
		jobStore, err = NewJobBoltStore(path.Join(c.RaftDataDir, "raft", "job.db"))
		if err != nil {
			logs.Fatal("NewJobBoltStore failed to create job store: err=%v", err)
		}
	}
	transport := raft.NewNetworkTransport(NewRaftStreamLayer(listener), raftNetworkTransportMaxPool,
		raftNetworkTransportTimeout, logs.GetRaftWriter())
	fsm := NewJobFSM(jobStore)
	underlay, err := raft.NewRaft(rc, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		logs.Fatal("NewRaftLayer failed to create raft instance: err=%v", err)
	}

	return &RaftLayer{
		bootstrap:     c.RaftBootstrap,
		underlay:      underlay,
		rc:            rc,
		snapshotStore: snapshotStore,
		stableStore:   stableStore,
		logStore:      logStore,
		transport:     transport,
	}
}

// Run starts raft layer.
func (l *RaftLayer) Run() {
	if l.bootstrap {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      l.rc.LocalID,
					Address: l.transport.LocalAddr(),
				},
			},
		}

		l.underlay.BootstrapCluster(configuration)
	}
}
