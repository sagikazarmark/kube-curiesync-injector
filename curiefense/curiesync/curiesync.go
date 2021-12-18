package curiesync

type RunMode string

func (rm RunMode) String() string {
	return string(rm)
}

const (
	SyncOnce      RunMode = "SYNC_ONCE"
	CopyBootstrap RunMode = "COPY_BOOTSTRAP"
	PeriodicSync  RunMode = "PERIODIC_SYNC"
)
