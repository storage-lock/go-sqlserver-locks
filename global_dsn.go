package sqlserver_locks

import (
	"context"
	"database/sql"
	sqlserver_storage "github.com/storage-lock/go-sqlserver-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
)

var dsnStorageLockFactoryBeanFactory *storage_lock_factory.StorageLockFactoryBeanFactory[string, *sql.DB] = storage_lock_factory.NewStorageLockFactoryBeanFactory[string, *sql.DB]()

// newSqlServerStorageByDsn 指向 sqlserver_storage.NewSqlServerStorage，包级变量以便单测注入构造结果，
// 覆盖 GetOrInit initFunc 的分支（无真库环境无法走通真实构造）。
var newSqlServerStorageByDsn = sqlserver_storage.NewSqlServerStorage

func NewSqlServerLockByDsn(ctx context.Context, dsn string, lockId string) (*storage_lock.StorageLock, error) {
	factory, err := GetSqlServerLockFactoryByDsn(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return factory.CreateLock(lockId)
}

func NewSqlServerLockByDsnWithOptions(ctx context.Context, uri string, options *storage_lock.StorageLockOptions) (*storage_lock.StorageLock, error) {
	factory, err := GetSqlServerLockFactoryByDsn(ctx, uri)
	if err != nil {
		return nil, err
	}
	return factory.CreateLockWithOptions(options)
}

func GetSqlServerLockFactoryByDsn(ctx context.Context, uri string) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
	return dsnStorageLockFactoryBeanFactory.GetOrInit(ctx, uri, func(ctx context.Context) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
		connectionManager := sqlserver_storage.NewSqlServerConnectionManagerFromDsn(uri)
		options := sqlserver_storage.NewSqlServerStorageOptions().SetConnectionManager(connectionManager)
		storage, err := newSqlServerStorageByDsn(ctx, options)
		if err != nil {
			return nil, err
		}
		factory := storage_lock_factory.NewStorageLockFactory(storage, options.ConnectionManager)
		return factory, nil
	})
}
