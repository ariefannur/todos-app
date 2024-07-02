package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Connect to PostgreSQL

	dbpool, err := pgxpool.NewWithConfig(context.Background(), Config())
	if err != nil {
		log.Fatal("Error while creating connection to the database!!")
	}
	defer dbpool.Close()

	connection, err := dbpool.Acquire(context.Background())
	if err != nil {
		log.Fatal("Error while acquiring connection from the database pool!! " + err.Error())
	}
	defer connection.Release()

	err = connection.Ping(context.Background())
	if err != nil {
		log.Fatal("Could not ping database")
	}

	fmt.Println("Connected to the database!!")

	createDBQuery := `
    CREATE TABLE IF NOT EXISTS todos (
        id SERIAL PRIMARY KEY,
        title VARCHAR(100) NOT NULL,
        description TEXT,
        completed BOOLEAN DEFAULT FALSE
    );
	`

	_, err = dbpool.Exec(context.Background(), createDBQuery)
	if err != nil {
		log.Fatalf("Failed to create database: %v\n", err)
	}

	fmt.Printf("Database '%s' created successfully.\n", "todos")

	app := fiber.New()

	// Define routes
	app.Get("/todos", getAllTodos(dbpool))
	app.Post("/todos", createTodo(dbpool))
	app.Get("/todos/:id", getTodo(dbpool))
	app.Put("/todos/:id", updateTodo(dbpool))
	app.Delete("/todos/:id", deleteTodo(dbpool))

	app.Listen(":8080")
}

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func getAllTodos(db *pgxpool.Pool) fiber.Handler {
	fn := func(c *fiber.Ctx) error {
		rows, err := db.Query(context.Background(), "SELECT id, title, description, completed FROM todos")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()

		var todos []Todo
		for rows.Next() {
			var todo Todo
			if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed); err != nil {
				return c.Status(500).SendString(err.Error())
			}
			todos = append(todos, todo)
		}
		return c.JSON(todos)
	}
	return fiber.Handler(fn)
}

func createTodo(db *pgxpool.Pool) fiber.Handler {
	fn := func(c *fiber.Ctx) error {
		todo := new(Todo)
		if err := c.BodyParser(todo); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		err := db.QueryRow(context.Background(),
			"INSERT INTO todos (title, description, completed) VALUES ($1, $2, $3) RETURNING id",
			todo.Title, todo.Description, todo.Completed).Scan(&todo.ID)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(todo)
	}
	return fiber.Handler(fn)
}

func getTodo(db *pgxpool.Pool) fiber.Handler {
	fn := func(c *fiber.Ctx) error {
		id := c.Params("id")
		var todo Todo
		err := db.QueryRow(context.Background(),
			"SELECT id, title, description, completed FROM todos WHERE id=$1", id).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed)
		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(404).SendString("Todo not found")
			}
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(todo)
	}

	return fiber.Handler(fn)

}

func updateTodo(db *pgxpool.Pool) fiber.Handler {
	fn := func(c *fiber.Ctx) error {
		id := c.Params("id")
		todo := new(Todo)
		if err := c.BodyParser(todo); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		commandTag, err := db.Exec(context.Background(),
			"UPDATE todos SET title=$1, description=$2, completed=$3 WHERE id=$4",
			todo.Title, todo.Description, todo.Completed, id)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		if commandTag.RowsAffected() == 0 {
			return c.Status(404).SendString("Todo not found")
		}
		return c.SendStatus(200)
	}
	return fiber.Handler(fn)
}

func deleteTodo(db *pgxpool.Pool) fiber.Handler {
	fn := func(c *fiber.Ctx) error {
		id := c.Params("id")
		commandTag, err := db.Exec(context.Background(),
			"DELETE FROM todos WHERE id=$1", id)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		if commandTag.RowsAffected() == 0 {
			return c.Status(404).SendString("Todo not found")
		}
		return c.SendStatus(200)
	}
	return fiber.Handler(fn)
}
