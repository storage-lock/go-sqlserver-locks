package sqlserver_locks

import (
	"context"
	storage_lock_test_helper "github.com/storage-lock/go-storage-lock-test-helper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewSqlServerLockByDsn(t *testing.T) {
	envName := "STORAGE_LOCK_SQLSERVER_DSN"
	sqlServerDsn := os.Getenv(envName)
	assert.NotEmpty(t, sqlServerDsn)

	factory, err := GetSqlServerLockFactoryByDsn(context.Background(), sqlServerDsn)
	assert.Nil(t, err)

	storage_lock_test_helper.PlayerNum = 10
	storage_lock_test_helper.EveryOnePlayTimes = 100
	storage_lock_test_helper.TestStorageLock(t, factory)
}
