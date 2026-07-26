package admin

import (
	"context"
	"fmt"
	"net/http"

	"mi-tramite-bolivia-backend/internal/auditoria"
	"mi-tramite-bolivia-backend/internal/db"
	"github.com/gin-gonic/gin"
)

// ListarCategoriasAdmin godoc
func ListarCategoriasAdmin(c *gin.Context) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT id, codigo, nombre, icono, orden, activa, padre_id
		FROM categoria
		ORDER BY orden ASC, nombre ASC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}
	defer rows.Close()

	var lista []map[string]interface{}
	for rows.Next() {
		var id int64
		var codigo, nombre string
		var icono *string
		var orden int
		var activa bool
		var padre_id *int64
		if err := rows.Scan(&id, &codigo, &nombre, &icono, &orden, &activa, &padre_id); err != nil {
			continue
		}
		lista = append(lista, map[string]interface{}{
			"id":       id,
			"codigo":   codigo,
			"nombre":   nombre,
			"icono":    icono,
			"orden":    orden,
			"activa":   activa,
			"padre_id": padre_id,
		})
	}
	c.JSON(http.StatusOK, gin.H{"datos": lista})
}

// CrearCategoria godoc
func CrearCategoria(c *gin.Context) {
	var req struct {
		Codigo  string  `json:"codigo" binding:"required"`
		Nombre  string  `json:"nombre" binding:"required"`
		Icono   *string `json:"icono"`
		Orden   int     `json:"orden"`
		PadreID *int64  `json:"padre_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	var id int64
	err := db.Pool.QueryRow(context.Background(), `
		INSERT INTO categoria (codigo, nombre, icono, orden, padre_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`, req.Codigo, req.Nombre, req.Icono, req.Orden, req.PadreID).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	idStr := fmt.Sprintf("%d", id)
	auditoria.Registrar(c.Request.Context(), actorID, "categoria.crear", "categoria", &idStr, nil, req, nil, nil)
	c.JSON(http.StatusCreated, gin.H{"message": "Categoría creada", "id": id})
}

// ActualizarCategoria godoc
func ActualizarCategoria(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Codigo  *string `json:"codigo"`
		Nombre  *string `json:"nombre"`
		Icono   *string `json:"icono"`
		Orden   *int    `json:"orden"`
		PadreID *int64  `json:"padre_id"`
		Activa  *bool   `json:"activa"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	_, err := db.Pool.Exec(context.Background(), `
		UPDATE categoria
		SET codigo = COALESCE($2, codigo),
		    nombre = COALESCE($3, nombre),
		    icono = COALESCE($4, icono),
		    orden = COALESCE($5, orden),
		    padre_id = COALESCE($6, padre_id),
		    activa = COALESCE($7, activa)
		WHERE id = $1
	`, id, req.Codigo, req.Nombre, req.Icono, req.Orden, req.PadreID, req.Activa)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	auditoria.Registrar(c.Request.Context(), actorID, "categoria.actualizar", "categoria", &id, nil, req, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Categoría actualizada"})
}

// EliminarCategoria godoc
func EliminarCategoria(c *gin.Context) {
	id := c.Param("id")
	
	_, err := db.Pool.Exec(context.Background(), `
		UPDATE categoria SET activa = false WHERE id = $1
	`, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	auditoria.Registrar(c.Request.Context(), actorID, "categoria.eliminar", "categoria", &id, nil, nil, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Categoría inactivada"})
}
