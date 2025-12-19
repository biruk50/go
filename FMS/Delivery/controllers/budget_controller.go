package controllers

import (
	"FMS/Domain"
	"FMS/Usecases"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BudgetController struct {
	BudgetUC Usecases.BudgetUsecase
}

func NewBudgetController(b Usecases.BudgetUsecase) *BudgetController {
	return &BudgetController{BudgetUC: b}
}

func (bc *BudgetController) GetAllBudgets(c *gin.Context) {
	budgets, err := bc.BudgetUC.GetAllBudgets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"budgets": budgets})
}

func (bc *BudgetController) GetBudgetByID(c *gin.Context) {
	id := c.Param("id")
	b, err := bc.BudgetUC.GetBudgetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"budget": b})
}

func (bc *BudgetController) GetBudgetSummary(c *gin.Context) {
	id := c.Param("id")
	summary, err := bc.BudgetUC.GetBudgetSummary(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"summary": summary})
}

func (bc *BudgetController) CreateBudget(c *gin.Context) {
	var payload Domain.Budget
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payload.CreatedAt = time.Now().UTC()
	created, err := bc.BudgetUC.CreateBudget(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"budget": created})
}

func (bc *BudgetController) UpdateBudget(c *gin.Context) {
	id := c.Param("id")
	var payload Domain.Budget
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := bc.BudgetUC.UpdateBudget(id, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (bc *BudgetController) ApproveBudget(c *gin.Context) {
	id := c.Param("id")
	if err := bc.BudgetUC.ApproveBudget(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "approved"})
}

func (bc *BudgetController) RejectBudget(c *gin.Context) {
	id := c.Param("id")
	if err := bc.BudgetUC.RejectBudget(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "rejected"})
}
