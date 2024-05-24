package todo

import (
	"context"
	"fmt"
	"log"
	"todo/ent"
	"todo/ent/todo"

	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3"
)

func Example_Todo() {
	// Create an ent.Client with in-memory SQLite database.
	client, err := ent.Open(
		dialect.SQLite,
		"file:ent?mode=memory&cache=shared&_fk=1",
	)
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	ctx := context.Background()
	// Run the automatic migration tool to create all schema resources.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	task1, err := client.Todo.Create().SetText("Add GraphQL Example").Save(ctx)
	if err != nil {
		log.Fatalf("failed creating a todo: %v", err)
	}
	fmt.Printf("%d: %q\n", task1.ID, task1.Text)

	task2, err := client.Todo.Create().SetText("Add Tracing Example").Save(ctx)
	if err != nil {
		log.Fatalf("failed creating a todo: %v", err)
	}
	fmt.Printf("%d: %q\n", task2.ID, task2.Text)

	if err := task2.Update().SetParent(task1).Exec(ctx); err != nil {
		log.Fatalf("failed connecting todo2 to its parent: %v", err)
	}

	items, err := client.Todo.Query().All(ctx)

	if err != nil {
		log.Fatalf("failed querying todos: %v", err)
	}

	for _, t := range items {
		fmt.Printf("%d: %q\n", t.ID, t.Text)
	}

	items, err = client.Todo.Query().Where(todo.HasParent()).All(ctx)
	if err != nil {
		log.Fatalf("failed querying todos: %v", err)
	}
	for _, t := range items {
		fmt.Printf("%d: %q\n", t.ID, t.Text)
	}

	items, err = client.Todo.Query().
		Where(
			todo.Not(
				todo.HasParent(),
			),
			todo.HasChildren(),
		).
		All(ctx)

	parent, err := client.Todo.Query().
		Where(todo.HasParent()).
		QueryParent().
		Only(ctx)
	if err != nil {
		log.Fatalf("failed querying todos: %v", err)
	}
	fmt.Printf("%d: %q\n", parent.ID, parent.Text)
}
