package controllers

import (
	"FMS/Domain"
	"FMS/Usecases"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CashRequestController struct {
	CashRequestUC Usecases.CashRequestUsecase
}

func NewCashRequestController(c Usecases.CashRequestUsecase) *CashRequestController {
	return &CashRequestController{CashRequestUC: c}
}

func (cc *CashRequestController) GetAllCashRequests(c *gin.Context) {
	list, err := cc.CashRequestUC.GetAllCashRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cash_requests": list})
}

func (cc *CashRequestController) GetCashRequest(c *gin.Context) {
	id := c.Param("id")
	r, err := cc.CashRequestUC.GetCashRequestByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cash_request": r})
}

func (cc *CashRequestController) CreateCashRequest(c *gin.Context) {
	var payload Domain.CashRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payload.CreatedAt = time.Now().UTC()
	created, err := cc.CashRequestUC.CreateCashRequest(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"cash_request": created})
}

func (cc *CashRequestController) ApproveCashRequest(c *gin.Context) {
	id := c.Param("id")
	if err := cc.CashRequestUC.ApproveCashRequest(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "approved"})
}

func (cc *CashRequestController) RejectCashRequest(c *gin.Context) {
	id := c.Param("id")
	if err := cc.CashRequestUC.RejectCashRequest(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "rejected"})
}

func (cc *CashRequestController) DisburseCashRequest(c *gin.Context) {
	id := c.Param("id")
	if err := cc.CashRequestUC.DisburseCashRequest(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "disbursed"})
}
