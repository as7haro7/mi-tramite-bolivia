package admin

import (
	"context"
	"fmt"
	"net/http"

	"mi-tramite-bolivia-backend/internal/auditoria"
	"mi-tramite-bolivia-backend/internal/db"

	"github.com/gin-gonic/gin"
)

// ListarFuentes godoc
// @Summary Listar fuentes de ingesta
// @Tags    Admin - Ingesta
// @Produce json
// @Security BearerAuth
// @Router  /api/v1/admin/fuentes [get]
func ListarFuentes(c *gin.Context) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT fi.id, fi.nombre, fi.tipo, fi.url, fi.frecuencia_cron,
		       fi.estado, fi.ultima_ejecucion_en, fi.proxima_ejecucion_en,
		       i.nombre AS institucion
		FROM fuente_ingesta fi
		LEFT JOIN institucion i ON i.id = fi.institucion_id
		ORDER BY fi.nombre
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}
	defer rows.Close()

	type FItem struct {
		db.FuenteIngesta
		Institucion *string `json:"institucion,omitempty"`
	}
	datos := make([]FItem, 0)
	for rows.Next() {
		var it FItem
		if err := rows.Scan(
			&it.ID, &it.Nombre, &it.Tipo, &it.URL, &it.FrecuenciaCron,
			&it.Estado, &it.UltimaEjecucionEn, &it.ProximaEjecucionEn,
			&it.Institucion,
		); err != nil {
			continue
		}
		datos = append(datos, it)
	}

	c.JSON(http.StatusOK, gin.H{"datos": datos})
}

// CrearFuente godoc
// @Summary Crear fuente de ingesta
// @Tags    Admin - Ingesta
// @Accept  json
// @Produce json
// @Security BearerAuth
// @Router  /api/v1/admin/fuentes [post]
func CrearFuente(c *gin.Context) {
	var req struct {
		InstitucionID  *int64  `json:"institucion_id"`
		Nombre         string  `json:"nombre" binding:"required"`
		Tipo           string  `json:"tipo" binding:"required"`
		URL            string  `json:"url" binding:"required"`
		FrecuenciaCron *string `json:"frecuencia_cron"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	var id int64
	err := db.Pool.QueryRow(context.Background(), `
		INSERT INTO fuente_ingesta (institucion_id, nombre, tipo, url, frecuencia_cron, estado)
		VALUES ($1, $2, $3, $4, $5, 'pausada') RETURNING id
	`, req.InstitucionID, req.Nombre, req.Tipo, req.URL, req.FrecuenciaCron).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	idStr := fmt.Sprintf("%d", id)
	auditoria.Registrar(c.Request.Context(), actorID, "fuente_ingesta.crear", "fuente_ingesta", &idStr, nil, req, nil, nil)
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// ActualizarFuente godoc
// @Summary Actualizar fuente de ingesta
// @Tags    Admin - Ingesta
// @Accept  json
// @Produce json
// @Security BearerAuth
// @Router  /api/v1/admin/fuentes/{id} [put]
func ActualizarFuente(c *gin.Context) {
	idStr := c.Param("id")
	var req struct {
		InstitucionID  *int64  `json:"institucion_id"`
		Nombre         *string `json:"nombre"`
		Tipo           *string `json:"tipo"`
		URL            *string `json:"url"`
		FrecuenciaCron *string `json:"frecuencia_cron"`
		Estado         *string `json:"estado"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	_, err := db.Pool.Exec(context.Background(), `
		UPDATE fuente_ingesta
		SET institucion_id = COALESCE($2, institucion_id),
		    nombre = COALESCE($3, nombre),
		    tipo = COALESCE($4, tipo),
		    url = COALESCE($5, url),
		    frecuencia_cron = COALESCE($6, frecuencia_cron),
		    estado = COALESCE($7, estado)
		WHERE id = $1
	`, idStr, req.InstitucionID, req.Nombre, req.Tipo, req.URL, req.FrecuenciaCron, req.Estado)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	auditoria.Registrar(c.Request.Context(), actorID, "fuente_ingesta.actualizar", "fuente_ingesta", &idStr, nil, req, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Fuente actualizada"})
}

// EliminarFuente godoc
// @Summary Eliminar fuente de ingesta
// @Tags    Admin - Ingesta
// @Produce json
// @Security BearerAuth
// @Router  /api/v1/admin/fuentes/{id} [delete]
func EliminarFuente(c *gin.Context) {
	idStr := c.Param("id")
	
	// Realmente borramos o la desactivamos? El esquema tiene `estado` que puede ser `error` o `pausada`.
	// Vamos a hacer un DELETE físico para simplificar si no está en uso fuerte, o un update a estado `pausada`. 
	// Dejémoslo en borrado físico, ya que es configuración.
	_, err := db.Pool.Exec(context.Background(), `
		DELETE FROM fuente_ingesta WHERE id = $1
	`, idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "No se puede eliminar la fuente porque tiene ingestas asociadas. Paúsela en su lugar."})
		return
	}

	actorID := actorIDFromCtx(c)
	auditoria.Registrar(c.Request.Context(), actorID, "fuente_ingesta.eliminar", "fuente_ingesta", &idStr, nil, nil, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Fuente eliminada"})
}

// ListarIngestas godoc
// @Summary Listar ejecuciones de ingesta
// @Tags    Admin - Ingesta
// @Produce json
// @Security BearerAuth
// @Router  /api/v1/admin/ingestas [get]
func ListarIngestas(c *gin.Context) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT ei.id, ei.fuente_ingesta_id, ei.estado,
		       ei.http_status, ei.registros_leidos, ei.candidatos_creados,
		       ei.mensaje_error, ei.iniciada_en, ei.finalizada_en,
		       fi.nombre AS fuente
		FROM ejecucion_ingesta ei
		JOIN fuente_ingesta fi ON fi.id = ei.fuente_ingesta_id
		ORDER BY ei.iniciada_en DESC
		LIMIT 100
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}
	defer rows.Close()

	type EjItem struct {
		db.EjecucionIngesta
		Fuente string `json:"fuente"`
	}
	datos := make([]EjItem, 0)
	for rows.Next() {
		var it EjItem
		if err := rows.Scan(
			&it.ID, &it.FuenteIngestaID, &it.Estado,
			&it.HTTPStatus, &it.RegistrosLeidos, &it.CandidatosCreados,
			&it.MensajeError, &it.IniciadaEn, &it.FinalizadaEn,
			&it.Fuente,
		); err != nil {
			continue
		}
		datos = append(datos, it)
	}

	c.JSON(http.StatusOK, gin.H{"datos": datos})
}

// ListarCandidatos godoc
// @Summary Listar candidatos de ingesta pendientes
// @Tags    Admin - Ingesta
// @Produce json
// @Security BearerAuth
// @Router  /api/v1/admin/candidatos [get]
func ListarCandidatos(c *gin.Context) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT ci.id, ci.ejecucion_id, ci.tramite_id_sugerido,
		       ci.datos_extraidos, ci.confianza, ci.estado,
		       ci.revisado_en
		FROM candidato_ingesta ci
		WHERE ci.estado = 'pendiente'
		ORDER BY ci.estado, ci.confianza DESC
		LIMIT 100
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}
	defer rows.Close()

	datos := make([]db.CandidatoIngesta, 0)
	for rows.Next() {
		var ci db.CandidatoIngesta
		if err := rows.Scan(
			&ci.ID, &ci.EjecucionID, &ci.TramiteIDSugerido,
			&ci.DatosExtraidos, &ci.Confianza, &ci.Estado,
			&ci.RevisadoEn,
		); err != nil {
			continue
		}
		datos = append(datos, ci)
	}

	c.JSON(http.StatusOK, gin.H{"datos": datos})
}

// ActualizarCandidato godoc
// @Summary Aceptar o rechazar candidato de ingesta
// @Tags    Admin - Ingesta
// @Accept  json
// @Produce json
// @Security BearerAuth
// @Param   id path string true "UUID del candidato"
// @Router  /api/v1/admin/candidatos/{id} [put]
func ActualizarCandidato(c *gin.Context) {
	candidatoID := c.Param("id")
	var req struct {
		Estado string `json:"estado" binding:"required"` // aceptado, rechazado, duplicado
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	_, err := db.Pool.Exec(context.Background(), `
		UPDATE candidato_ingesta
		SET estado = $2, revisado_por = $3, revisado_en = NOW()
		WHERE id = $1 AND estado = 'pendiente'
	`, candidatoID, req.Estado, actorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}

	auditoria.Registrar(c.Request.Context(), actorID, "candidato."+req.Estado, "candidato_ingesta", &candidatoID, nil, req, nil, nil)

	// Si se aceptó, podría generarse un borrador automáticamente en una versión futura.
	// Por ahora solo cambia el estado; el editor crea la versión manualmente.
	c.JSON(http.StatusOK, gin.H{"message": "Candidato actualizado a: " + req.Estado})
}
