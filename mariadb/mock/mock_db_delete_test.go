package mock

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

// Test_DB_Delete mocks database and deletes operation with a mock database.
func Test_DB_Delete(t *testing.T) {
	// Set up mock database and SQL mock
	mockDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	// Defer closing the mock database connection
	defer func() {
		_ = mockDB.Close()
		// require.NoError(t, err)
	}()

	// Expect a query with specific arguments and return the mock response
	sqlMock.ExpectExec("DELETE FROM users WHERE id = (.+)").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Execute the query
	var result sql.Result
	result, err = mockDB.Exec("DELETE FROM users WHERE id = ?", 1)
	require.NoError(t, err)

	// Assert the expected affect number
	affect, err := result.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), affect)
}
