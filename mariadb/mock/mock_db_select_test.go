package mock

import (
	"fmt"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Alice").
		AddRow(2, "Bob")

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id = (.+)").
		WithArgs(1).
		WillReturnRows(rows)

	var id int
	var name string
	err = db.QueryRow("SELECT id, name FROM users WHERE id = ?", 1).Scan(&id, &name)
	if err != nil {
		t.Fatalf("Error querying database: %v", err)
	}

	fmt.Printf("User %d: %s\n", id, name)
}

type User struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
