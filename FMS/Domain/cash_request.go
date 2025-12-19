package Domain

import(
	"time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type CashRequest struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description,omitempty" json:"description"`
	Amount      float64            `bson:"amount" json:"amount"`
	BudgetID    primitive.ObjectID `bson:"budget_id,omitempty" json:"budget_id,omitempty"`
	Requester   string             `bson:"requester,omitempty" json:"requester,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Status      string             `bson:"status,omitempty" json:"status"`
}
