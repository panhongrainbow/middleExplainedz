# The Three Features of DB SQL

DB SQL has three major features:

1. `Implementation of Connection Pool`
   DB SQL internally implements a connection pool, which allows for efficient management and reuse of database connections.
   This helps optimize connection establishment and `reduces the overhead of creating new connections for each request`.
2. `Concurrency Safety with Locking` 
   DB SQL ensures the safety of concurrent access through the use of locks.
   `This allows multiple goroutines to safely access the database concurrently` without conflicts or data corruption.
3. `Automatic Data Type Conversion`
   DB SQL provides `automatic data type conversion`.
   This means that it can handle the conversion of data types between the database and the application code seamlessly, making it easier to work with different data types and reducing the need for manual type conversions in the code.

## Feature 1: Implementation of Connection Pool

Retrieve old connection

```go
func (db *DB) conn(ctx context.Context, strategy connReuseStrategy) (*driverConn, error) {
    // ...
	last := len(db.freeConn) - 1 // <<<<< reuse the old db connection
	if strategy == cachedOrNewConn && last >= 0 {
		conn := db.freeConn[last] // <<<<< reuse the old db connection
        db.freeConn = db.freeConn[:last] // <<<<< reuse the old db connection (the smallest idle time)
		conn.inUse = true
		if conn.expired(lifetime) {
			db.maxLifetimeClosed++
			db.mu.Unlock()
			conn.Close()
			return nil, driver.ErrBadConn
		}
```

Store old connection

```go
func (db *DB) putConnDBLocked(dc *driverConn, err error) bool {
	if db.closed {
		return false
	}
	if db.maxOpen > 0 && db.numOpen > db.maxOpen {
		return false
	}
	if c := len(db.connRequests); c > 0 {
        // ... (no db sql requests)
	} else if err == nil && !db.closed {
        // ... adds the current connection to the pool
		if db.maxIdleConnsLocked() > len(db.freeConn) {
			db.freeConn = append(db.freeConn, dc) // adds the connection to the pool
			db.startCleanerLocked()
			return true
		}
		db.maxIdleClosed++
	}
	return false
}
```

For `performance reasons`, since `the connection with the minimum idle time has been selected` from the idle connections, it can be `ensured that the returned connection is valid and not expired`.

(由池中选出最小的闲置连线中，确保可以立即使用未过期的连线)

Because connections must have an expiration time set, but retrieving expired connections from the pool would decrease performance, such a design is necessary.

(一切都是为了效能)

## Feature 2: Concurrency Safety with Locking

Immediately support high concurrency

```go
func (db *DB) conn(ctx context.Context, strategy connReuseStrategy) (*driverConn, error) {
    db.mu.Lock() // <<<<< Only one goroutine at a time can retrieve a connection from the pool
    if db.closed {
        db.mu.Unlock()
        return nil, errDBClosed
    }
    // ...
    db.mu.Unlock() // <<<<< Only one goroutine at a time can retrieve a connection from the pool.
}
```

## Feature 3: Automatic Data Type Conversion

The simplest example

```go
package mock

import (
	"fmt"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func TestGetUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Alice").
		AddRow(2, "Bob")

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id = (.+)").
		WithArgs(1).
		WillReturnRows(rows)

	var id int
	var name string
	_ = db.QueryRow("SELECT id, name FROM users WHERE id = ?", 1).Scan(&id, &name)
    // <<<<< Using the Scan function in Golang to perform Type Conversion
    // 使用 SCAN 乙式进行 Type Conversion
    // The final value of id is 1 and the value of name is Alice.
}
```

## Additional information

FieldsByTraversal is an important function used to help SQLX perform type conversion

```go
package mock

import (
	"errors"
	"github.com/jmoiron/sqlx/reflectx"
	"reflect"
)

func FieldsByTraversal(v reflect.Value, traversals [][]int, values []interface{}, ptrs bool) error {
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return errors.New("argument not a struct")
	}

	for i, traversal := range traversals {
		if len(traversal) == 0 {
			values[i] = new(interface{})
			continue
		}
		f := reflectx.FieldByIndexes(v, traversal)
		if ptrs {
			values[i] = f.Addr().Interface()
		} else {
			values[i] = f.Interface()
		}
	}
	return nil
}
```

Using a function called `FieldsByTraversal`, `the person struct and values are synchronized`.

When `db sql` uses `Scan function` to update the values slice, `the values of the person struct will also be modified`.

(把 person struct 跟 values 进行同步)

```go
import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestFieldsByTraversal2(t *testing.T) {
	person := Person{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
	}

	traversals := [][]int{
		{0}, // Name
		{1}, // Age
		{2}, // Email
	}

	values := make([]interface{}, len(traversals))

	_ = FieldsByTraversal(reflect.ValueOf(&person), traversals, values, true)

	require.Equal(t, "Alice", *values[0].(*string))
	require.Equal(t, 30, *values[1].(*int))
	require.Equal(t, "alice@example.com", *values[2].(*string))
}

```

