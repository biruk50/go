package Domain

import(
	"time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Budget represents an allocated budget for a period or department
type Budget struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description,omitempty" json:"description"`
	Amount      float64            `bson:"amount" json:"amount"`
	Remaining   float64            `bson:"remaining" json:"remaining"`
	Department  string             `bson:"department,omitempty" json:"department,omitempty"`
	DueDate     time.Time          `bson:"due_date,omitempty" json:"due_date"`
	Status      string             `bson:"status,omitempty" json:"status"`
	CreatedBy   string             `bson:"created_by,omitempty" json:"created_by,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}
