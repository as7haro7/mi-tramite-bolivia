package admin

import (
	"context"
	"fmt"
	"net/http"

	"mi-tramite-bolivia-backend/internal/auditoria"
	"mi-tramite-bolivia-backend/internal/db"
	"github.com/gin-gonic/gin"
)

// ListarEtiquetasAdmin godoc
func ListarEtiquetasAdmin(c *gin.Context) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT id, slug, nombre
		FROM etiqueta
		ORDER BY nombre ASC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}
	defer rows.Close()

	var lista []map[string]interface{}
	for rows.Next() {
		var id int64
		var slug, nombre string
		if err := rows.Scan(&id, &slug, &nombre); err != nil {
			continue
		}
		lista = append(lista, map[string]interface{}{
			"id":     id,
			"slug":   slug,
			"nombre": nombre,
		})
	}
	c.JSON(http.StatusOK, gin.H{"datos": lista})
}

// CrearEtiqueta godoc
func CrearEtiqueta(c *gin.Context) {
	var req struct {
		Slug   string `json:"slug" binding:"required"`
		Nombre string `json:"nombre" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	var id int64
	err := db.Pool.QueryRow(context.Background(), `
		INSERT INTO etiqueta (slug, nombre)
		VALUES ($1, $2) RETURNING id
	`, req.Slug, req.Nombre).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	idStr := fmt.Sprintf("%d", id)
	auditoria.Registrar(c.Request.Context(), actorID, "etiqueta.crear", "etiqueta", &idStr, nil, req, nil, nil)
	c.JSON(http.StatusCreated, gin.H{"message": "Etiqueta creada", "id": id})
}

// ActualizarEtiqueta godoc
func ActualizarEtiqueta(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Slug   *string `json:"slug"`
		Nombre *string `json:"nombre"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	_, err := db.Pool.Exec(context.Background(), `
		UPDATE etiqueta
		SET slug = COALESCE($2, slug),
		    nombre = COALESCE($3, nombre)
		WHERE id = $1
	`, id, req.Slug, req.Nombre)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	auditoria.Registrar(c.Request.Context(), actorID, "etiqueta.actualizar", "etiqueta", &id, nil, req, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Etiqueta actualizada"})
}

// EliminarEtiqueta godoc
// @Summary Eliminar etiqueta (hard delete porque no tiene estado)
func EliminarEtiqueta(c *gin.Context) {
	id := c.Param("id")
	
	_, err := db.Pool.Exec(context.Background(), `
		DELETE FROM etiqueta WHERE id = $1
	`, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	auditoria.Registrar(c.Request.Context(), actorID, "etiqueta.eliminar", "etiqueta", &id, nil, nil, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Etiqueta eliminada"})
}
