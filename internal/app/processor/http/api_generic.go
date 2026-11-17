package rprocessor

import (
	"net/http"

	rhandler "github.com/AlexanderTal/order-service/internal/app/handler/http"
	"github.com/gin-gonic/gin"
)

func vGenericRegHealthCheck(r *gin.Engine, h rhandler.Health) {
	r.GET("/health", h.LastCheck)
}

func handleNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
}
