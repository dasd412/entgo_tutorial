package entgo_tutorial

import (
	"context"
	"entdemo/ent"
	"entdemo/ent/car"
	"entdemo/ent/group"
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

func QueryCarsUsers(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)

	if err != nil {
		return fmt.Errorf("failed querying cars: %w", err)
	}

	//역방향 엣지 질의
	for _, c := range cars {
		owner, err := c.QueryOwner().Only(ctx)

		if err != nil {
			return fmt.Errorf("failed querying cars: %w", err)
		}

		log.Printf("car %s owner %s\n", c.Model, owner.Name)
	}

	return nil
}

func CreateGraph(ctx context.Context, client *ent.Client) error {
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("Ariel").
		Save(ctx)

	if err != nil {
		return err
	}

	neta, err := client.User.
		Create().
		SetAge(28).
		SetName("Neta").
		Save(ctx)

	if err != nil {
		return err
	}

	err = client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		SetOwner(a8m).
		Exec(ctx)

	if err != nil {
		return err
	}

	err = client.Car.
		Create().
		SetModel("Mazda").
		SetRegisteredAt(time.Now()).
		SetOwner(a8m).
		Exec(ctx)

	if err != nil {
		return err
	}

	err = client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		SetOwner(neta).
		Exec(ctx)

	if err != nil {
		return err
	}

	err = client.Group.
		Create().
		SetName("GitLab").
		AddUsers(neta, a8m).
		Exec(ctx)

	if err != nil {
		return err
	}

	err = client.Group.
		Create().
		SetName("GitHub").
		AddUsers(a8m).
		Exec(ctx)

	if err != nil {
		return err
	}

	log.Println("The graph was created")

	return nil
}

func QueryGithub(ctx context.Context, client *ent.Client) error {
	cars, err := client.Group.
		Query().
		Where(group.Name("Github")).
		QueryUsers().
		QueryCars().
		All(ctx)

	if err != nil {
		return fmt.Errorf("failed getting cars: %w", err)
	}

	log.Println("cars returned : ", cars)

	return nil
}

func QueryArielCars(ctx context.Context, client *ent.Client) error {
	a8m := client.User.
		Query().
		Where(user.HasCars(), user.Name("Ariel")).
		OnlyX(ctx)

	cars, err := a8m.QueryGroups().
		QueryUsers().
		QueryCars().
		Where(car.Not(car.Model("Mazda"))).
		All(ctx)

	if err != nil {
		return fmt.Errorf("failed getting cars: %w", err)
	}

	log.Println("cars returned : ", cars)

	return nil
}

func QueryGroupWithUsers(ctx context.Context, client *ent.Client) error {
	groups, err := client.Group.
		Query().
		Where(group.HasUsers()).
		All(ctx)

	if err != nil {
		return fmt.Errorf("failed getting groups: %w", err)
	}

	log.Println("groups returned : ", groups)

	return nil
}
