package Domain

import(
	"time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description,omitempty" json:"description"`
	DueDate     time.Time          `bson:"due_date,omitempty" json:"due_date"`
	Status      string             `bson:"status,omitempty" json:"status"`
}

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Username     string             `bson:"username" json:"username"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Role         string             `bson:"role" json:"role"` // "admin" or "user"
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}
