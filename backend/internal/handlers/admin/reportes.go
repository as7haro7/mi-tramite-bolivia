package admin

import (
	"context"
	"fmt"
	"net/http"

	"mi-tramite-bolivia-backend/internal/auditoria"
	"mi-tramite-bolivia-backend/internal/db"

	"github.com/gin-gonic/gin"
)

// ListarReportesAdmin godoc
// @Summary Listar reportes ciudadanos
// @Tags    Admin - Reportes
// @Produce json
// @Security BearerAuth
// @Param   estado query string false "Filtrar por estado (nuevo, en_revision, resuelto, descartado)"
// @Router  /api/v1/admin/reportes [get]
func ListarReportesAdmin(c *gin.Context) {
	estado := c.Query("estado")
	args := []interface{}{}
	where := ""
	if estado != "" {
		where = "WHERE rc.estado = $1"
		args = append(args, estado)
	}

	rows, err := db.Pool.Query(context.Background(), `
		SELECT rc.id, rc.tramite_id, rc.oficina_id, rc.tipo,
		       rc.descripcion, rc.estado, rc.asignado_a,
		       rc.creado_en, rc.resuelto_en,
		       t.slug AS tramite_slug
		FROM reporte_ciudadano rc
		LEFT JOIN tramite t ON t.id = rc.tramite_id
		`+where+`
		ORDER BY rc.creado_en DESC
		LIMIT 200
	`, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}
	defer rows.Close()

	type RItem struct {
		db.ReporteCiudadano
		TramiteSlug *string `json:"tramite_slug,omitempty"`
	}
	datos := make([]RItem, 0)
	for rows.Next() {
		var r RItem
		if err := rows.Scan(
			&r.ID, &r.TramiteID, &r.OficinaID, &r.Tipo,
			&r.Descripcion, &r.Estado, &r.AsignadoA,
			&r.CreadoEn, &r.ResueltoEn,
			&r.TramiteSlug,
		); err != nil {
			continue
		}
		datos = append(datos, r)
	}

	c.JSON(http.StatusOK, gin.H{"datos": datos})
}

// ActualizarReporte godoc
// @Summary Actualizar estado de reporte ciudadano
// @Tags    Admin - Reportes
// @Accept  json
// @Produce json
// @Security BearerAuth
// @Param   id path string true "UUID del reporte"
// @Router  /api/v1/admin/reportes/{id} [put]
func ActualizarReporte(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Estado    string  `json:"estado" binding:"required"`
		AsignadoA *int64  `json:"asignado_a"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	resueltoEn := "NULL"
	if req.Estado == "resuelto" {
		resueltoEn = "NOW()"
	}

	_, err := db.Pool.Exec(context.Background(), fmt.Sprintf(`
		UPDATE reporte_ciudadano
		SET estado = $2,
		    asignado_a = COALESCE($3, asignado_a),
		    resuelto_en = CASE WHEN $2 = 'resuelto' THEN %s ELSE resuelto_en END
		WHERE id = $1
	`, resueltoEn), id, req.Estado, req.AsignadoA)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	auditoria.Registrar(c.Request.Context(), actorID, "reporte.actualizar", "reporte_ciudadano", &id, nil, req, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Reporte actualizado"})
}
