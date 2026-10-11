package sqlserver_locks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalDsnLockFactoryBeanFactory_InitiallyInitialized(t *testing.T) {
	assert.NotNil(t, dsnStorageLockFactoryBeanFactory)
}

func TestNewSqlServerLockByDsn_InvalidDsn(t *testing.T) {
	_, err := NewSqlServerLockByDsn(context.Background(), "invalid-dsn", "test-lock")
	assert.Error(t, err)
}

func TestGetSqlServerLockFactoryByDsn_InvalidDsn(t *testing.T) {
	_, err := GetSqlServerLockFactoryByDsn(context.Background(), "invalid-dsn")
	assert.Error(t, err)
}

func TestNewSqlServerLockByDsnWithOptions_InvalidDsn(t *testing.T) {
	_, err := NewSqlServerLockByDsnWithOptions(context.Background(), "invalid-dsn", nil)
	assert.Error(t, err)
}

