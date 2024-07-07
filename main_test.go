package main

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTodosTable(t *testing.T) {
	// Create a new mock pool
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockPool.Close()

	// Expect the table creation query
	mockPool.ExpectExec("CREATE TABLE IF NOT EXISTS todos").WillReturnResult(pgxmock.NewResult("CREATE", 1))

	// Call the function to create the table`
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS todos (
        id SERIAL PRIMARY KEY,
        title VARCHAR(100) NOT NULL,
        description TEXT,
        completed BOOLEAN DEFAULT FALSE
    );
	`
	// Call the function to create the table
	_, err = mockPool.Exec(context.Background(), createTableQuery)
	assert.NoError(t, err)

	// Ensure all expectations were met
	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)

}
