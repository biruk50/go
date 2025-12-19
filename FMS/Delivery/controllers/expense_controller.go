package controllers

import (
	"FMS/Domain"
	"FMS/Usecases"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ExpenseController struct {
	ExpenseUC Usecases.ExpenseUsecase
}

func NewExpenseController(e Usecases.ExpenseUsecase) *ExpenseController {
	return &ExpenseController{ExpenseUC: e}
}

func (ec *ExpenseController) GetAllExpenses(c *gin.Context) {
	list, err := ec.ExpenseUC.GetAllExpenses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"expenses": list})
}

func (ec *ExpenseController) GetExpense(c *gin.Context) {
	id := c.Param("id")
	e, err := ec.ExpenseUC.GetExpenseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"expense": e})
}

func (ec *ExpenseController) GetExpenseSummary(c *gin.Context) {
	id := c.Param("id")
	// simple summary: return expense details for now
	e, err := ec.ExpenseUC.GetExpenseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"summary": e})
}

func (ec *ExpenseController) CreateExpense(c *gin.Context) {
	var payload Domain.Expense
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payload.CreatedAt = time.Now().UTC()
	created, err := ec.ExpenseUC.CreateExpense(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"expense": created})
}

func (ec *ExpenseController) CreateExpenseReceipt(c *gin.Context) {
	id := c.Param("id")
	var payload struct {
		ReceiptURL string `json:"receipt_url"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil || payload.ReceiptURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "receipt_url required"})
		return
	}
	if err := ec.ExpenseUC.AttachReceipt(id, payload.ReceiptURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "receipt attached"})
}

func (ec *ExpenseController) VerifyExpense(c *gin.Context) {
	id := c.Param("id")
	if err := ec.ExpenseUC.VerifyExpense(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "verified"})
}
