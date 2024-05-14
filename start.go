package entgo_tutorial

import (
	"context"
	"entdemo/ent"
	"entdemo/ent/user"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	client, err := ent.Open("mysql", "<user>:<pass>@tcp(<host>:<port>)/<database>?parseTime=True")

	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}

	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	user, err := client.User.Create().SetAge(30).SetName("a8m").Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}

	log.Println("user was created", user)

	return user, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.Name("a8m")).
		Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}

	log.Println("user returned : ", u)

	return u, nil
}
