package state

import (
	"context"
	"math/rand"
	"time"

	"github.com/hashicorp/terraform/states/statemgr"
)

var rngSource *rand.Rand

func init() {
	rngSource = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// State is a deprecated alias for statemgr.Full.
type State = statemgr.Full

// StateReader is a deprecated alias for statemgr.Reader
type StateReader = statemgr.Reader

// StateWriter is a deprecated alias for stagemgr.Writer
type StateWriter = statemgr.Writer

// StateRefresher is a deprecated alias for statemgr.Refresher
type StateRefresher = statemgr.Refresher

// StatePersister is a deprecated alias for statemgr.Persister
type StatePersister = statemgr.Persister

// Locker is a deprecated alias for statemgr.Locker
type Locker = statemgr.Locker

// test hook to verify that LockWithContext has attempted a lock
var postLockHook func()

// LockWithContext is a deprecated alias for statemgr.LockWithContext.
func LockWithContext(ctx context.Context, s State, info *LockInfo) (string, error) {
	return statemgr.LockWithContext(ctx, s, info)
}

// NewLockInfo is a deprecated alias for statemgr.NewLockInfo.
func NewLockInfo() *LockInfo {
	return statemgr.NewLockInfo()
}

// LockInfo is a deprecated alias for statemgr.LockInfo
type LockInfo = statemgr.LockInfo

// LockError is a deprecated alias for statemgr.LockError
type LockError = statemgr.LockError
