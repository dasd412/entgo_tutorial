package entgo_tutorial

import (
	"context"
	"entdemo/ent"
	"entdemo/ent/car"
	"entdemo/ent/user"
	"fmt"
	"log"
	"time"

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
	u, err := client.User.Create().SetAge(30).SetName("a8m").Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}

	log.Println("user was created", u)

	return u, nil
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

func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
	tesla, err := client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to creating ca: %w", err)
	}

	log.Println("car was created", tesla)

	ford, err := client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}

	log.Println("car was created", ford)

	// 유저 1명을 만든 후, 2개의 차와 연관 관계를 맺는다.
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		AddCars(tesla, ford).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed creating user:  %w", err)
	}

	log.Println("user was created", a8m)
	return a8m, nil
}

func QueryCars(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)

	if err != nil {
		return fmt.Errorf("failed querying cars: %w", err)
	}

	log.Println("cars returned : ", cars)

	//특정 차량만 필터링하기
	ford, err := a8m.QueryCars().
		Where(car.Model("Ford")).
		Only(ctx)

	if err != nil {
		return fmt.Errorf("failed querying cars: %w", err)
	}

	log.Println(ford)

	return nil
}
