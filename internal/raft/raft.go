package raft

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

// Layer represents crond raft consensus.
type Layer struct {
	bootstrap bool
	underlay  *raft.Raft

	rc            *raft.Config
	snapshotStore raft.SnapshotStore
	stableStore   raft.StableStore
	logStore      raft.LogStore
	transport     raft.Transport
}

// NewLayer creates crond raft Layer.
func NewLayer(c *LayerConfig, listener net.Listener) *Layer {
	rc := raft.DefaultConfig()
	rc.LogOutput = logs.GetRaftWriter()
	rc.LocalID = raft.ServerID(c.RaftNode)

	var err error

	var snapshotStore raft.SnapshotStore
	var stableStore raft.StableStore
	var logStore raft.LogStore

	if !c.RaftProduction {
		snapshotStore = raft.NewInmemSnapshotStore()

		memStore := raft.NewInmemStore()
		stableStore = memStore
		logStore = memStore
	} else {
		snapshotStore, err = raft.NewFileSnapshotStore(path.Join(c.RaftDataDir, "raft"), raftFileSnapshotStoreRetain, logs.GetRaftWriter())
		if err != nil {
			logs.Fatal("NewLayer failed to create snapshotStore: err=%v", err)
		}

		boltStore, err := raftboltdb.NewBoltStore(path.Join(c.RaftDataDir, "raft", "raft.db"))
		if err != nil {
			logs.Fatal("NewLayer failed to create stable store: err=%v", err)
		}

		stableStore = boltStore

		logStore, err = raft.NewLogCache(raftMaxLogCacheSize, boltStore)
		if err != nil {
			logs.Fatal("NewLayer failed to create log store: err=%v", err)
		}
	}
	transport := raft.NewNetworkTransport(NewStreamLayer(listener), raftNetworkTransportMaxPool,
		raftNetworkTransportTimeout, logs.GetRaftWriter())

	underlay, err := raft.NewRaft(rc, nil, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		logs.Fatal("NewLayer failed to create raft instance: err=%v", err)
	}

	return &Layer{
		bootstrap:     c.RaftBootstrap,
		underlay:      underlay,
		rc:            rc,
		snapshotStore: snapshotStore,
		stableStore:   stableStore,
		logStore:      logStore,
		transport:     transport,
	}
}

// Run starts raft.
func (l *Layer) Run() {
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
