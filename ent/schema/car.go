package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Car holds the schema definition for the Car entity.
type Car struct {
	ent.Schema
}

// Fields of the Car.
func (Car) Fields() []ent.Field {
	return []ent.Field{
		field.String("model"),
		field.Time("registered_at"),
	}
}

// Edges of the Car.
func (Car) Edges() []ent.Edge {
	// 차에서 사용자 방향으로 역방향 엣지 설정하기
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("cars").
			Unique(), //차 한대에 오직 사용자 하나만 매핑되도록 유니크 엣지로 세팅
	}
}
