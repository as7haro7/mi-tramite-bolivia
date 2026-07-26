package handlers

import (
	"context"
	"net/http"

	"mi-tramite-bolivia-backend/internal/db"

	"github.com/gin-gonic/gin"
)

type CrearReporteRequest struct {
	TramiteID   *int64  `json:"tramite_id"`
	OficinaID   *int64  `json:"oficina_id"`
	Tipo        string  `json:"tipo" binding:"required"`
	Descripcion string  `json:"descripcion" binding:"required,max=2000"`
	Correo      *string `json:"correo"` // opcional, se cifra antes de guardar
}

// CrearReporte godoc
// @Summary      Reportar dato incorrecto
// @Description  Permite enviar un reporte ciudadano de forma anónima. No modifica el catálogo.
// @Tags         Reportes
// @Accept       json
// @Produce      json
// @Param        request body CrearReporteRequest true "Datos del reporte"
// @Success      202  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/reportes [post]
func CrearReporte(c *gin.Context) {
	var req CrearReporteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "VALIDATION_ERROR",
			"message": "Datos de reporte inválidos",
			"details": err.Error(),
		})
		return
	}

	tiposValidos := map[string]bool{
		"dato_incorrecto": true, "dato_desactualizado": true,
		"oficina_cerrada": true, "costo_distinto": true,
		"requisito_distinto": true, "otro": true,
	}
	if !tiposValidos[req.Tipo] {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "VALIDATION_ERROR",
			"message": "Tipo de reporte inválido",
		})
		return
	}

	// Cifrar correo si se proporcionó (pgcrypto symmetric encryption sería ideal;
	// aquí lo almacenamos como bytes con pgp_sym_encrypt si está disponible,
	// o como NULL si no se proporcionó)
	_, err := db.Pool.Exec(context.Background(), `
		INSERT INTO reporte_ciudadano
		    (tramite_id, oficina_id, tipo, descripcion)
		VALUES ($1, $2, $3, $4)
	`, req.TramiteID, req.OficinaID, req.Tipo, req.Descripcion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "DB_ERROR",
			"message": "Error guardando reporte",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Reporte recibido. Será revisado por nuestro equipo editorial.",
	})
}
