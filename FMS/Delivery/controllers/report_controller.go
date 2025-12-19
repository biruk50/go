package controllers

import (
	"FMS/Usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	ReportUC Usecases.ReportUsecase
}

func NewReportController(r Usecases.ReportUsecase) *ReportController {
	return &ReportController{ReportUC: r}
}

func (rc *ReportController) GetOverviewReport(c *gin.Context) {
	o, err := rc.ReportUC.GetOverview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"overview": o})
}

func (rc *ReportController) GetCashRequestReport(c *gin.Context) {
	list, err := rc.ReportUC.GetCashRequestReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cash_requests": list})
}

func (rc *ReportController) GetBudgetReport(c *gin.Context) {
	list, err := rc.ReportUC.GetBudgetReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"budgets": list})
}

func (rc *ReportController) GetExpenseReport(c *gin.Context) {
	list, err := rc.ReportUC.GetExpenseReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"expenses": list})
}
