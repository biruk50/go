package Domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Expense struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description,omitempty" json:"description"`
	Amount      float64            `bson:"amount" json:"amount"`
	ReceiptURL  string             `bson:"receipt_url,omitempty" json:"receipt_url,omitempty"`
	BudgetID    primitive.ObjectID `bson:"budget_id,omitempty" json:"budget_id,omitempty"`
	CreatedBy   string             `bson:"created_by,omitempty" json:"created_by,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	DueDate     time.Time          `bson:"due_date,omitempty" json:"due_date"`
	Status      string             `bson:"status,omitempty" json:"status"`
}
