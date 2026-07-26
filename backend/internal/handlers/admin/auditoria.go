package admin

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"mi-tramite-bolivia-backend/internal/db"

	"github.com/gin-gonic/gin"
)

// ListarAuditoria godoc
// @Summary Listar eventos de auditoría
// @Description Filtrable por actor, acción, entidad y fecha. Solo lectura; no se puede modificar desde el panel.
// @Tags    Admin - Auditoría
// @Produce json
// @Security BearerAuth
// @Param   actor     query int    false "ID del actor"
// @Param   accion    query string false "Acción (ej: version.publicar)"
// @Param   entidad   query string false "Tipo de entidad (ej: tramite_version)"
// @Param   desde     query string false "Fecha desde (ISO 8601)"
// @Param   hasta     query string false "Fecha hasta (ISO 8601)"
// @Param   limite    query int    false "Límite (máx 100)"
// @Router  /api/v1/admin/auditoria [get]
func ListarAuditoria(c *gin.Context) {
	args := []interface{}{}
	idx := 1
	where := []string{}

	if actor := c.Query("actor"); actor != "" {
		where = append(where, "ea.actor_id = $"+strconv.Itoa(idx))
		args = append(args, actor)
		idx++
	}
	if accion := c.Query("accion"); accion != "" {
		where = append(where, "ea.accion ILIKE '%'||$"+strconv.Itoa(idx)+`||'%'`)
		args = append(args, accion)
		idx++
	}
	if entidad := c.Query("entidad"); entidad != "" {
		where = append(where, "ea.entidad_tipo = $"+strconv.Itoa(idx))
		args = append(args, entidad)
		idx++
	}
	if desde := c.Query("desde"); desde != "" {
		where = append(where, "ea.ocurrido_en >= $"+strconv.Itoa(idx))
		args = append(args, desde)
		idx++
	}
	if hasta := c.Query("hasta"); hasta != "" {
		where = append(where, "ea.ocurrido_en <= $"+strconv.Itoa(idx))
		args = append(args, hasta)
		idx++
	}

	limite := 50
	if l, err := strconv.Atoi(c.Query("limite")); err == nil && l > 0 && l <= 100 {
		limite = l
	}
	args = append(args, limite)

	whereSQL := ""
	if len(where) > 0 {
		whereSQL = "WHERE " + strings.Join(where, " AND ")
	}

	rows, err := db.Pool.Query(context.Background(), `
		SELECT ea.id, ea.actor_id, u.correo AS actor_correo,
		       ea.accion, ea.entidad_tipo, ea.entidad_id,
		       ea.antes, ea.despues, ea.ocurrido_en
		FROM evento_auditoria ea
		LEFT JOIN usuario_admin u ON u.id = ea.actor_id
		`+whereSQL+`
		ORDER BY ea.ocurrido_en DESC
		LIMIT $`+strconv.Itoa(idx),
		args...,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}
	defer rows.Close()

	type EvItem struct {
		db.EventoAuditoria
		ActorCorreo *string `json:"actor_correo,omitempty"`
	}
	datos := make([]EvItem, 0)
	for rows.Next() {
		var ev EvItem
		if err := rows.Scan(
			&ev.ID, &ev.ActorID, &ev.ActorCorreo,
			&ev.Accion, &ev.EntidadTipo, &ev.EntidadID,
			&ev.Antes, &ev.Despues, &ev.OcurridoEn,
		); err != nil {
			continue
		}
		datos = append(datos, ev)
	}

	c.JSON(http.StatusOK, gin.H{"datos": datos})
}
