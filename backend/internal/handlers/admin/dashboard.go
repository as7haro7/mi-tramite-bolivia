package admin

import (
	"context"
	"net/http"

	"mi-tramite-bolivia-backend/internal/db"

	"github.com/gin-gonic/gin"
)

// Dashboard godoc
// @Summary Dashboard operativo del panel administrativo
// @Description Devuelve las métricas clave para la bandeja del editor: versiones pendientes, trámites vencidos, embeddings fallidos, reportes nuevos, etc.
// @Tags    Admin - Dashboard
// @Produce json
// @Security BearerAuth
// @Router  /api/v1/admin/dashboard [get]
func Dashboard(c *gin.Context) {
	type Metrics struct {
		VersionesEnRevision      int `json:"versiones_en_revision"`
		TramitesRevisionVencida  int `json:"tramites_revision_vencida"`
		EmbeddingsFallidos       int `json:"embeddings_fallidos"`
		FuentesConError          int `json:"fuentes_con_error"`
		ReportesCiudadanosNuevos int `json:"reportes_ciudadanos_nuevos"`
		CandidatosPendientes     int `json:"candidatos_pendientes"`
	}

	ctx := context.Background()
	var m Metrics

	_ = db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM tramite_version WHERE estado_editorial = 'en_revision'
	`).Scan(&m.VersionesEnRevision)

	_ = db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM tramite_version
		WHERE estado_editorial = 'publicada'
		  AND proxima_revision_en IS NOT NULL
		  AND proxima_revision_en < NOW()
	`).Scan(&m.TramitesRevisionVencida)

	_ = db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM trabajo_embedding WHERE estado = 'fallido'
	`).Scan(&m.EmbeddingsFallidos)

	_ = db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM fuente_ingesta WHERE estado = 'error'
	`).Scan(&m.FuentesConError)

	_ = db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM reporte_ciudadano WHERE estado = 'nuevo'
	`).Scan(&m.ReportesCiudadanosNuevos)

	_ = db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM candidato_ingesta WHERE estado = 'pendiente'
	`).Scan(&m.CandidatosPendientes)

	// Publicaciones por institución
	type InstPub struct {
		Institucion   string `json:"institucion"`
		Publicaciones int    `json:"publicaciones"`
	}
	rows, _ := db.Pool.Query(ctx, `
		SELECT i.nombre, COUNT(*) AS publicaciones
		FROM tramite_version v
		JOIN tramite t ON t.id = v.tramite_id
		JOIN institucion i ON i.id = t.institucion_id
		WHERE v.estado_editorial = 'publicada'
		GROUP BY i.nombre
		ORDER BY publicaciones DESC
	`)
	var porInstitucion []InstPub
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var it InstPub
			if err := rows.Scan(&it.Institucion, &it.Publicaciones); err == nil {
				porInstitucion = append(porInstitucion, it)
			}
		}
	}

	if porInstitucion == nil {
		porInstitucion = []InstPub{}
	}

	c.JSON(http.StatusOK, gin.H{
		"metricas":       m,
		"por_institucion": porInstitucion,
	})
}
