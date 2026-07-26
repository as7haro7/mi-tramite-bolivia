package handlers

import (
	"context"
	"net/http"

	"mi-tramite-bolivia-backend/internal/db"

	"github.com/gin-gonic/gin"
)

// HealthLive godoc
// @Summary      Liveness check
// @Description  Devuelve 200 siempre que el proceso esté corriendo.
// @Tags         Health
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Router       /health/live [get]
func HealthLive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// HealthReady godoc
// @Summary      Readiness check
// @Description  Devuelve 200 si el servidor puede atender peticiones (BD disponible).
// @Tags         Health
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      503  {object} map[string]interface{}
// @Router       /health/ready [get]
func HealthReady(c *gin.Context) {
	if err := db.Pool.Ping(context.Background()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"error":  "database unavailable",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
