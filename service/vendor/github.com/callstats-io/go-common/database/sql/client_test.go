package database_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	database "github.com/callstats-io/go-common/database/sql"
	"github.com/callstats-io/go-common/testutil"
	"github.com/callstats-io/go-common/vault"
)

func TestNewVaultSQLClient(t *testing.T) {
	fakeVaultReader := &vault.MockReader{}
	fakeVaultReader.MockSecret(map[string]interface{}{})

	testutil.MustBeNil(t, testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
		if _, err := database.NewVaultSQLClient(ctx, fakeVaultReader); err != nil {
			t.Errorf("Expected NewVaultSQLClient not to return an error on initialize")
		}
	}))
}
func TestDB(t *testing.T) {
	mockValidSecret := func(c *vault.MockReader) {
		c.MockSecret(map[string]interface{}{
			"username": "user",
			"password": "pass",
		})
	}

	t.Run("with canceled call context", func(t *testing.T) {
		t.Parallel()
		fakeVaultReader := &vault.MockReader{}
		mockValidSecret(fakeVaultReader)
		fakeVaultReader.BlockReads()

		testutil.MustBeNil(t, testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
			client, err := database.NewVaultSQLClient(ctx, fakeVaultReader)
			testutil.MustBeNil(t, err)

			// create a closed context
			callCtx, callCtxCancel := context.WithCancel(ctx)
			callCtxCancel()

			if _, err := client.DB(callCtx); err != callCtx.Err() {
				fmt.Println(callCtx.Err(), err)
				t.Errorf("Expected canceled context to return error")
			}
		}))
	})

	t.Run("with closed client", func(t *testing.T) {
		t.Parallel()
		fakeVaultReader := &vault.MockReader{}
		mockValidSecret(fakeVaultReader)
		fakeVaultReader.BlockReads()

		testutil.MustBeNil(t, testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
			client, err := database.NewVaultSQLClient(ctx, fakeVaultReader)
			testutil.MustBeNil(t, err)
			// close the client and expect DB to fail
			client.Close()

			testutil.MustBeNil(t, testutil.WithDeadlineContext(time.Millisecond, func(callCtx context.Context) {
				if _, err := client.DB(callCtx); err != database.ErrClosed {
					t.Errorf("Expected close to close the client")
				}
			}))
		}))
	})

	t.Run("with error from vault", func(t *testing.T) {
		t.Parallel()
		fakeVaultReader := &vault.MockReader{}
		expErr := errors.New("EXPERR")
		fakeVaultReader.MockError(expErr)

		testutil.MustBeNil(t, testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
			client, err := database.NewVaultSQLClient(ctx, fakeVaultReader)
			testutil.MustBeNil(t, err)

			testutil.MustBeNil(t, testutil.WithDeadlineContext(time.Millisecond, func(callCtx context.Context) {
				if _, err := client.DB(callCtx); err != expErr {
					t.Errorf("Expected vault error to propagate up")
				}
			}))
		}))
	})

	t.Run("with error from sql driver", func(t *testing.T) {
		t.Parallel()

		fakeVaultReader := &vault.MockReader{}
		mockValidSecret(fakeVaultReader)

		testutil.MustBeNil(t, testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
			client, err := database.NewVaultSQLClient(ctx, fakeVaultReader, database.OptionSQLDriver("test-driver-not-registered"))
			testutil.MustBeNil(t, err)

			testutil.MustBeNil(t, testutil.WithDeadlineContext(time.Millisecond, func(callCtx context.Context) {
				if _, err := client.DB(callCtx); err.Error() != "sql: unknown driver \"test-driver-not-registered\" (forgotten import?)" {
					t.Errorf("Expected sql driver error to propagate up")
				}
			}))
		}))
	})

	t.Run("with error from db connection", func(t *testing.T) {
		t.Parallel()
		driverName := "test-sql-driver-error"
		fakeDriver := &fakeSQLDriver{}
		sql.Register(driverName, fakeDriver)
		expErr := errors.New("EXPERR")
		fakeDriver.MockError(expErr)

		fakeVaultReader := &vault.MockReader{}
		mockValidSecret(fakeVaultReader)

		testutil.MustBeNil(t, testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
			client, err := database.NewVaultSQLClient(ctx, fakeVaultReader, database.OptionSQLDriver(driverName))
			testutil.MustBeNil(t, err)

			testutil.MustBeNil(t, testutil.WithDeadlineContext(time.Millisecond, func(callCtx context.Context) {
				if _, err := client.DB(callCtx); err != expErr {
					t.Errorf("Expected db connect error to propagate up")
				}
			}))
		}))
	})

	t.Run("with valid defaulted setup", func(t *testing.T) {
		t.Parallel()
		fakeVaultReader := &vault.MockReader{}
		mockValidSecret(fakeVaultReader)

		driverName := "postgres"
		fakeDriver := &fakeSQLDriver{}
		sql.Register(driverName, fakeDriver)

		testutil.MustBeNil(t, testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
			client, err := database.NewVaultSQLClient(ctx, fakeVaultReader)
			testutil.MustBeNil(t, err)

			testutil.MustBeNil(t, testutil.WithDeadlineContext(time.Millisecond, func(callCtx context.Context) {
				db, err := client.DB(callCtx)
				testutil.MustBeNil(t, err)
				if db == nil {
					t.Errorf("Expected to get a fake postgres connection")
				}
				testutil.MustBeNil(t, db.Ping())
			}))
		}))
	})

	t.Run("with multiple calls", func(t *testing.T) {
		t.Parallel()
		fakeVaultReader := &vault.MockReader{}
		mockValidSecret(fakeVaultReader)

		driverName := "test-valid-multi-call-driver"
		fakeDriver := &fakeSQLDriver{}
		sql.Register(driverName, fakeDriver)

		testutil.MustBeNil(t, testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
			client, err := database.NewVaultSQLClient(ctx, fakeVaultReader, database.OptionSQLDriver(driverName))
			testutil.MustBeNil(t, err)

			testutil.MustBeNil(t, testutil.WithDeadlineContext(time.Millisecond, func(callCtx context.Context) {
				_, err := client.DB(callCtx)
				testutil.MustBeNil(t, err)
				if fakeDriver.connCalls != 1 {
					t.Errorf("Expected to increase conn calls count in fake driver")
				}

				// reset and call again, expect the client not to call the actual driver and instead use the cached connection
				fakeDriver.Reset()
				_, err = client.DB(callCtx)
				testutil.MustBeNil(t, err)

				if fakeDriver.connCalls != 0 {
					t.Errorf("Expected to get same postgres connection")
				}
			}))
		}))
	})
}
