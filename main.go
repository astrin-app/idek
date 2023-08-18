package main

import (
	"astrin/main/db"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/mustache/v2"
)

type Todo struct {
	ID        int
	Task      string `json:"task"`
	Completed bool
}

func main() {
	dbUrl := os.Getenv("DBURL")
	myDB, err := db.CreateDB(dbUrl)
	if err != nil {
		fmt.Print(err)
	}
	engine := mustache.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Get("/", func(c *fiber.Ctx) error {
		rows, err := myDB.Query("SELECT * FROM todos")
		if err != nil {
			fmt.Println(err)
			return c.Status(500).SendString("Error fetching todos")
		}
		defer rows.Close()

		todos := []Todo{} // Create a slice to hold the todos

		// Loop through the rows and scan them into the Todo struct
		for rows.Next() {
			var todo Todo
			err := rows.Scan(&todo.ID, &todo.Task, &todo.Completed)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(todo)
			todos = append(todos, todo)
		}
		return c.Render("index", fiber.Map{
			"Todos": todos, // Pass the todos to the template
		})
	})
	app.Post("/create-todo", func(c *fiber.Ctx) error {
		task := c.FormValue("task")
		res, err := myDB.Exec(fmt.Sprintf(`INSERT INTO TODOS (TASK, COMPLETED)
		VALUES ('%s', FALSE);`, task))
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(400)
		}
		lastId, _ := res.LastInsertId()
		fmt.Println(lastId)
		return c.Render("todo", fiber.Map{
			"Task": task,
			"ID":   lastId,
		})
	})
	app.Delete("/delete/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		myDB.Exec(fmt.Sprintf("delete from todos where id = %s", id))
		return c.Render("delete", fiber.Map{})
	})
	app.Listen(":6969")
}
