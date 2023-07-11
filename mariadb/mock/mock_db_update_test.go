package mock

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

// Test_DB_Update mocks database and updates operation with a mock database.
func Test_DB_Update(t *testing.T) {
	// Set up mock database and SQL mock
	mockDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	// Defer closing the mock database connection
	defer func() {
		_ = mockDB.Close()
		// require.NoError(t, err)
	}()

	// Set up mock response
	sqlMock.ExpectExec("UPDATE users SET name = (.+) WHERE id = (.+)").
		WithArgs("Bob", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Execute the query
	var result sql.Result
	result, err = mockDB.Exec("UPDATE users SET name = ? WHERE id = ?", "Bob", 1)
	require.NoError(t, err)

	// Assert the expected affect number
	affect, err := result.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), affect)
}
