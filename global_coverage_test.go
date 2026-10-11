package sqlserver_locks

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/storage-lock/go-storage"
	sqlserver_storage "github.com/storage-lock/go-sqlserver-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"github.com/stretchr/testify/assert"
)

func TestNewSqlServerLockBySqlDb(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer mockDB.Close()

	// Init 建表
	mock.ExpectExec("IF NOT EXISTS").
		WillReturnResult(sqlmock.NewResult(0, 0))

	ctx := context.Background()
	factory, err := GetSqlServerLockFactoryBySqlDb(ctx, mockDB)
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	// 缓存命中
	factory2, err := GetSqlServerLockFactoryBySqlDb(ctx, mockDB)
	assert.Nil(t, err)
	assert.Equal(t, factory, factory2)

	lock, err := NewSqlServerLockBySqlDb(ctx, mockDB, "test-sqlserver-lock-sqldb")
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	options := storage_lock.NewStorageLockOptionsWithLockId("test-sqlserver-lock-sqldb-opt")
	lock2, err := NewSqlServerLockBySqlDbWithOptions(ctx, mockDB, options)
	assert.Nil(t, err)
	assert.NotNil(t, lock2)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestNewSqlServerLockBySqlDb_Errors(t *testing.T) {
	old := newSqlServerStorage
	newSqlServerStorage = func(ctx context.Context, options *sqlserver_storage.SqlServerStorageOptions) (*sqlserver_storage.SqlServerStorage, error) {
		return nil, errors.New("boom")
	}
	defer func() { newSqlServerStorage = old }()

	mockDB, _, err := sqlmock.New()
	assert.Nil(t, err)
	defer mockDB.Close()

	ctx := context.Background()
	_, err = GetSqlServerLockFactoryBySqlDb(ctx, mockDB)
	assert.NotNil(t, err)
	_, err = NewSqlServerLockBySqlDb(ctx, mockDB, "x")
	assert.NotNil(t, err)
	_, err = NewSqlServerLockBySqlDbWithOptions(ctx, mockDB, storage_lock.NewStorageLockOptionsWithLockId("x"))
	assert.NotNil(t, err)
}

func TestNewSqlServerLockByDsn_Success(t *testing.T) {
	// 无真库环境：注入构造成功，用 sqlmock 撑起真实的建表 Init
	old := newSqlServerStorageByDsn
	newSqlServerStorageByDsn = func(ctx context.Context, options *sqlserver_storage.SqlServerStorageOptions) (*sqlserver_storage.SqlServerStorage, error) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			return nil, err
		}
		mock.ExpectExec("IF NOT EXISTS").
			WillReturnResult(sqlmock.NewResult(0, 0))
		opts := sqlserver_storage.NewSqlServerStorageOptions().
			SetConnectionManager(storage.NewFixedSqlDBConnectionManager(mockDB))
		return sqlserver_storage.NewSqlServerStorage(ctx, opts)
	}
	defer func() { newSqlServerStorageByDsn = old }()

	ctx := context.Background()
	dsn := "mock-sqlserver-dsn-success"
	factory, err := GetSqlServerLockFactoryByDsn(ctx, dsn)
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	// 缓存命中
	factory2, err := GetSqlServerLockFactoryByDsn(ctx, dsn)
	assert.Nil(t, err)
	assert.Equal(t, factory, factory2)

	lock, err := NewSqlServerLockByDsn(ctx, dsn, "test-sqlserver-lock-dsn")
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	options := storage_lock.NewStorageLockOptionsWithLockId("test-sqlserver-lock-dsn-opt")
	lock2, err := NewSqlServerLockByDsnWithOptions(ctx, dsn, options)
	assert.Nil(t, err)
	assert.NotNil(t, lock2)
}

func TestNewSqlServerLockByDsn_InjectError(t *testing.T) {
	old := newSqlServerStorageByDsn
	newSqlServerStorageByDsn = func(ctx context.Context, options *sqlserver_storage.SqlServerStorageOptions) (*sqlserver_storage.SqlServerStorage, error) {
		return nil, errors.New("boom")
	}
	defer func() { newSqlServerStorageByDsn = old }()

	ctx := context.Background()
	dsn := "mock-sqlserver-dsn-error"
	_, err := GetSqlServerLockFactoryByDsn(ctx, dsn)
	assert.NotNil(t, err)
	_, err = NewSqlServerLockByDsn(ctx, dsn, "x")
	assert.NotNil(t, err)
	_, err = NewSqlServerLockByDsnWithOptions(ctx, dsn, storage_lock.NewStorageLockOptionsWithLockId("x"))
	assert.NotNil(t, err)
}
