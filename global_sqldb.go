package sqlserver_locks

import (
	"context"
	"database/sql"
	sqlserver_storage "github.com/storage-lock/go-sqlserver-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
)

var sqlDbStorageLockFactoryBeanFactory *storage_lock_factory.StorageLockFactoryBeanFactory[*sql.DB, *sql.DB] = storage_lock_factory.NewStorageLockFactoryBeanFactory[*sql.DB, *sql.DB]()

func NewSqlServerLockBySqlDb(ctx context.Context, db *sql.DB, lockId string) (*storage_lock.StorageLock, error) {
	factory, err := GetSqlServerLockFactoryBySqlDb(ctx, db)
	if err != nil {
		return nil, err
	}
	return factory.CreateLock(lockId)
}

func NewSqlServerLockBySqlDbWithOptions(ctx context.Context, db *sql.DB, options *storage_lock.StorageLockOptions) (*storage_lock.StorageLock, error) {
	factory, err := GetSqlServerLockFactoryBySqlDb(ctx, db)
	if err != nil {
		return nil, err
	}
	return factory.CreateLockWithOptions(options)
}

func GetSqlServerLockFactoryBySqlDb(ctx context.Context, db *sql.DB) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
	return sqlDbStorageLockFactoryBeanFactory.GetOrInit(ctx, db, func(ctx context.Context) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
		connectionManager := sqlserver_storage.NewSqlServerConnectionManagerFromSqlDb(db)
		options := sqlserver_storage.NewSqlServerStorageOptions().SetConnectionManager(connectionManager)
		storage, err := sqlserver_storage.NewSqlServerStorage(ctx, options)
		if err != nil {
			return nil, err
		}
		factory := storage_lock_factory.NewStorageLockFactory(storage, options.ConnectionManager)
		return factory, nil
	})
}
