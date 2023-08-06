package postgresql_locks

import (
	storage_lock "github.com/storage-lock/go-storage-lock"
	"sync"
)

var GlobalSqlServerLockFactory *SqlServerLockFactory
var globalSqlServerLockFactoryOnce sync.Once
var globalSqlServerLockFactoryErr error

func InitGlobalSqlServerLockFactory(uri string) error {
	factory, err := NewSqlServerLockFactory(uri)
	if err != nil {
		return err
	}
	GlobalSqlServerLockFactory = factory
	return nil
}

func NewSqlServerLock(uri string, lockId string) (*storage_lock.StorageLock, error) {
	globalSqlServerLockFactoryOnce.Do(func() {
		globalSqlServerLockFactoryErr = InitGlobalSqlServerLockFactory(uri)
	})
	if globalSqlServerLockFactoryErr != nil {
		return nil, globalSqlServerLockFactoryErr
	}
	return GlobalSqlServerLockFactory.CreateLock(lockId)
}
