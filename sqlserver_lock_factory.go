package postgresql_locks

import (
	"context"
	"database/sql"
	sqlserver_storage "github.com/storage-lock/go-sqlserver-storage"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
)

type SqlServerLockFactory struct {
	*storage_lock_factory.StorageLockFactory[*sql.DB]
}

func NewSqlServerLockFactory(dsn string) (*SqlServerLockFactory, error) {

	connectionManager := sqlserver_storage.NewSqlServerStorageConnectionGetterFromDSN(dsn)
	options := sqlserver_storage.NewSqlServerStorageOptions().SetConnectionManage(connectionManager)
	mongoStorage, err := sqlserver_storage.NewSqlServerStorage(context.Background(), options)
	if err != nil {
		return nil, err
	}
	factory := storage_lock_factory.NewStorageLockFactory[*sql.DB](mongoStorage, nil)
	return &SqlServerLockFactory{
		StorageLockFactory: factory,
	}, nil
}
