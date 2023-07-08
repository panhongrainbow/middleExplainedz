package mock

import (
	"github.com/stretchr/testify/require"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

// Test_DB_Select_One mocks database and executes a query to select a single row with specific arguments.
func Test_DB_Select_One(t *testing.T) {
	// Set up mock database and SQL mock
	mockDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	// Defer closing the mock database connection
	defer func() {
		_ = mockDB.Close()
		// require.NoError(t, err)
	}()

	// Set up mock response
	mockResponse := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Alice")

	// Expect a query with specific arguments and return the mock response
	sqlMock.ExpectQuery("SELECT (.+) FROM users WHERE id = (.+)").
		WithArgs(1).
		WillReturnRows(mockResponse)

	// Execute the query and scan the results into variables
	var id int
	var name string
	err = mockDB.QueryRow("SELECT id, name FROM users WHERE id = ?", 1).Scan(&id, &name)
	require.NoError(t, err)

	// Assert the expected values
	require.Equal(t, 1, id)
	require.Equal(t, "Alice", name)
}

// Test_DB_Select mocks database and executes a query to select multiple rows, iterating over the results.
func Test_DB_Select(t *testing.T) {
	// Set up mock database and SQL mock
	mockDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	// Defer closing the mock database connection
	defer func() {
		_ = mockDB.Close()
		// require.NoError(t, err)
	}()

	// Set up mock response
	mockResponse := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Alice").
		AddRow(2, "Bob")

	// Expect a query with specific arguments and return the mock response
	sqlMock.ExpectQuery("SELECT (.+) FROM users").
		WillReturnRows(mockResponse)

	// Execute the query and scan the results into variables
	type User struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}
	var users []User
	rows, err := mockDB.Query("SELECT id, name FROM users")
	require.NoError(t, err)

	// Iterate over the rows and scan the values into user objects
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Name)
		require.NoError(t, err)
		users = append(users, user)
	}

	// Assert the expected number of users
	require.Equal(t, 2, len(users))

	// Assert the values of the first user
	require.Equal(t, 1, users[0].ID)
	require.Equal(t, "Alice", users[0].Name)

	// Assert the values of the second user
	require.Equal(t, 2, users[1].ID)
	require.Equal(t, "Bob", users[1].Name)
}
