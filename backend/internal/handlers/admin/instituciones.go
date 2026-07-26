package admin

import (
	"context"
	"fmt"
	"net/http"

	"mi-tramite-bolivia-backend/internal/auditoria"
	"mi-tramite-bolivia-backend/internal/db"

	"github.com/gin-gonic/gin"
)

// ─── Instituciones ─────────────────────────────────────────────────────────────

// ListarInstitucionesAdmin godoc
// @Summary Listar instituciones (admin)
// @Tags    Admin - Instituciones
// @Produce json
// @Security BearerAuth
// @Router  /api/v1/admin/instituciones [get]
func ListarInstitucionesAdmin(c *gin.Context) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT i.id, i.codigo, i.nombre, i.sigla, i.tipo, i.sitio_web,
		       i.estado, i.creado_en, i.actualizado_en,
		       am.url_publica AS logo_url
		FROM institucion i
		LEFT JOIN archivo_multimedia am ON am.id = i.logo_archivo_id
		    AND am.variante = 'original' AND am.estado = 'activo'
		ORDER BY i.nombre
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}
	defer rows.Close()

	type InstAdmin struct {
		db.Institucion
		LogoURL *string `json:"logo_url,omitempty"`
	}

	datos := make([]InstAdmin, 0)
	for rows.Next() {
		var it InstAdmin
		if err := rows.Scan(
			&it.ID, &it.Codigo, &it.Nombre, &it.Sigla, &it.Tipo, &it.SitioWeb,
			&it.Estado, &it.CreadoEn, &it.ActualizadoEn, &it.LogoURL,
		); err != nil {
			continue
		}
		datos = append(datos, it)
	}

	c.JSON(http.StatusOK, gin.H{"datos": datos})
}

// CrearInstitucionAdmin godoc
// @Summary Crear institución
// @Tags    Admin - Instituciones
// @Accept  json
// @Produce json
// @Security BearerAuth
// @Router  /api/v1/admin/instituciones [post]
func CrearInstitucionAdmin(c *gin.Context) {
	var req struct {
		Codigo  string  `json:"codigo" binding:"required"`
		Nombre  string  `json:"nombre" binding:"required"`
		Sigla   string  `json:"sigla" binding:"required"`
		Tipo    string  `json:"tipo"`
		SitioWeb *string `json:"sitio_web"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}
	if req.Tipo == "" {
		req.Tipo = "publica"
	}

	var id int64
	err := db.Pool.QueryRow(context.Background(), `
		INSERT INTO institucion (codigo, nombre, sigla, tipo, sitio_web)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`, req.Codigo, req.Nombre, req.Sigla, req.Tipo, req.SitioWeb).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error creando institución: " + err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	idStr := fmt.Sprintf("%d", id)
	auditoria.Registrar(c.Request.Context(), actorID, "institucion.crear", "institucion", &idStr, nil, req, nil, nil)
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// ActualizarInstitucionAdmin godoc
// @Summary Actualizar institución
// @Tags    Admin - Instituciones
// @Accept  json
// @Produce json
// @Security BearerAuth
// @Param   id path int true "ID de la institución"
// @Router  /api/v1/admin/instituciones/{id} [put]
func ActualizarInstitucionAdmin(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Nombre   *string `json:"nombre"`
		Sigla    *string `json:"sigla"`
		SitioWeb *string `json:"sitio_web"`
		Estado   *string `json:"estado"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	_, err := db.Pool.Exec(context.Background(), `
		UPDATE institucion
		SET nombre   = COALESCE($2, nombre),
		    sigla    = COALESCE($3, sigla),
		    sitio_web = COALESCE($4, sitio_web),
		    estado   = COALESCE($5, estado)
		WHERE id = $1
	`, id, req.Nombre, req.Sigla, req.SitioWeb, req.Estado)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	auditoria.Registrar(c.Request.Context(), actorID, "institucion.actualizar", "institucion", &id, nil, req, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Institución actualizada"})
}

// EliminarInstitucionAdmin godoc
// @Summary Eliminar institución
// @Tags    Admin - Instituciones
// @Produce json
// @Security BearerAuth
// @Param   id path int true "ID de la institución"
// @Router  /api/v1/admin/instituciones/{id} [delete]
func EliminarInstitucionAdmin(c *gin.Context) {
	id := c.Param("id")
	
	_, err := db.Pool.Exec(context.Background(), `
		UPDATE institucion SET estado = 'inactiva' WHERE id = $1
	`, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	auditoria.Registrar(c.Request.Context(), actorID, "institucion.eliminar", "institucion", &id, nil, nil, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Institución inactivada"})
}

// ─── Oficinas ─────────────────────────────────────────────────────────────────

// ListarOficinasAdmin godoc
// @Summary Listar oficinas (admin)
// @Tags    Admin - Oficinas
// @Produce json
// @Security BearerAuth
// @Router  /api/v1/admin/oficinas [get]
func ListarOficinasAdmin(c *gin.Context) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT o.id, o.codigo, o.nombre, o.tipo, o.direccion,
		       o.latitud, o.longitud, o.requiere_cita, o.estado,
		       o.verificado_en, i.nombre AS institucion,
		       m.nombre AS municipio
		FROM oficina o
		JOIN institucion i ON i.id = o.institucion_id
		LEFT JOIN municipio m ON m.id = o.municipio_id
		ORDER BY o.nombre
		LIMIT 500
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}
	defer rows.Close()

	type OfiItem struct {
		db.OficinaSummary
		Codigo       string  `json:"codigo"`
		Institucion  string  `json:"institucion"`
		VerificadoEn *string `json:"verificado_en,omitempty"`
	}

	datos := make([]db.OficinaSummary, 0)
	for rows.Next() {
		var o db.OficinaSummary
		var codigo, institucion string
		var verificado *string
		if err := rows.Scan(
			&o.ID, &codigo, &o.Nombre, &o.Tipo, &o.Direccion,
			&o.Latitud, &o.Longitud, &o.RequiereCita, &o.Estado,
			&verificado, &institucion, &o.Municipio,
		); err != nil {
			continue
		}
		_ = codigo
		_ = institucion
		_ = verificado
		datos = append(datos, o)
	}

	c.JSON(http.StatusOK, gin.H{"datos": datos})
}

// CrearOficinaAdmin godoc
// @Summary Crear oficina
// @Tags    Admin - Oficinas
// @Accept  json
// @Produce json
// @Security BearerAuth
// @Router  /api/v1/admin/oficinas [post]
func CrearOficinaAdmin(c *gin.Context) {
	var req struct {
		InstitucionID int64    `json:"institucion_id" binding:"required"`
		MunicipioID   *int64   `json:"municipio_id"`
		Codigo        string   `json:"codigo" binding:"required"`
		Nombre        string   `json:"nombre" binding:"required"`
		Tipo          string   `json:"tipo"`
		Direccion     *string  `json:"direccion"`
		Latitud       *float64 `json:"latitud"`
		Longitud      *float64 `json:"longitud"`
		RequiereCita  bool     `json:"requiere_cita"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}
	if req.Tipo == "" {
		req.Tipo = "oficina"
	}

	// Validar coordenadas si se proporcionan
	if req.Latitud != nil && (*req.Latitud < -90 || *req.Latitud > 90) {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": "Latitud fuera de rango"})
		return
	}
	if req.Longitud != nil && (*req.Longitud < -180 || *req.Longitud > 180) {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": "Longitud fuera de rango"})
		return
	}

	var id int64
	err := db.Pool.QueryRow(context.Background(), `
		INSERT INTO oficina
		    (institucion_id, municipio_id, codigo, nombre, tipo,
		     direccion, latitud, longitud, requiere_cita)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, req.InstitucionID, req.MunicipioID, req.Codigo, req.Nombre, req.Tipo,
		req.Direccion, req.Latitud, req.Longitud, req.RequiereCita,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error creando oficina: " + err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	idStr := fmt.Sprintf("%d", id)
	auditoria.Registrar(c.Request.Context(), actorID, "oficina.crear", "oficina", &idStr, nil, req, nil, nil)
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// ActualizarOficinaAdmin godoc
// @Summary Actualizar oficina
// @Tags    Admin - Oficinas
// @Accept  json
// @Produce json
// @Security BearerAuth
// @Param   id path int true "ID de la oficina"
// @Router  /api/v1/admin/oficinas/{id} [put]
func ActualizarOficinaAdmin(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Nombre       *string  `json:"nombre"`
		Direccion    *string  `json:"direccion"`
		Latitud      *float64 `json:"latitud"`
		Longitud     *float64 `json:"longitud"`
		Estado       *string  `json:"estado"`
		RequiereCita *bool    `json:"requiere_cita"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	_, err := db.Pool.Exec(context.Background(), `
		UPDATE oficina
		SET nombre        = COALESCE($2, nombre),
		    direccion     = COALESCE($3, direccion),
		    latitud       = COALESCE($4, latitud),
		    longitud      = COALESCE($5, longitud),
		    estado        = COALESCE($6, estado),
		    requiere_cita = COALESCE($7, requiere_cita),
		    verificado_en = NOW()
		WHERE id = $1
	`, id, req.Nombre, req.Direccion, req.Latitud, req.Longitud, req.Estado, req.RequiereCita)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	auditoria.Registrar(c.Request.Context(), actorID, "oficina.actualizar", "oficina", &id, nil, req, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Oficina actualizada"})
}

// EliminarOficinaAdmin godoc
// @Summary Eliminar oficina
// @Tags    Admin - Oficinas
// @Produce json
// @Security BearerAuth
// @Param   id path int true "ID de la oficina"
// @Router  /api/v1/admin/oficinas/{id} [delete]
func EliminarOficinaAdmin(c *gin.Context) {
	id := c.Param("id")
	
	_, err := db.Pool.Exec(context.Background(), `
		UPDATE oficina SET estado = 'inactiva' WHERE id = $1
	`, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	auditoria.Registrar(c.Request.Context(), actorID, "oficina.eliminar", "oficina", &id, nil, nil, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Oficina inactivada"})
}
