package handlers

import (
	"context"
	"net/http"
	"time"

	"mi-tramite-bolivia-backend/internal/db"

	"github.com/gin-gonic/gin"
)

// GetTramiteBySlug godoc
// @Summary      Ficha completa de un trámite
// @Description  Devuelve la versión publicada vigente de un trámite por slug, incluyendo modalidades, requisitos, pasos, costos, resultados y fuentes.
// @Tags         Catálogo
// @Produce      json
// @Param        slug path string true "Slug del trámite"
// @Success      200  {object} db.TramiteFicha
// @Failure      404  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/tramites/{slug} [get]
func GetTramiteBySlug(c *gin.Context) {
	slug := c.Param("slug")

	// Leer tramite + versión publicada
	var ficha db.TramiteFicha
	var inst db.InstitucionResumen
	var versionID int64
	err := db.Pool.QueryRow(context.Background(), `
		SELECT
		    t.id, t.slug, v.titulo, v.resumen, v.descripcion,
		    v.publico_objetivo, v.resultado_esperado, v.advertencias,
		    v.plazo_texto, v.vigencia_resultado_texto,
		    v.requiere_cita, v.url_inicio,
		    v.verificado_en, v.proxima_revision_en,
		    i.codigo, i.nombre, i.sigla, i.sitio_web,
		    c.nombre AS categoria,
		    v.id AS version_id
		FROM tramite t
		JOIN institucion i ON i.id = t.institucion_id
		LEFT JOIN categoria c ON c.id = t.categoria_id
		JOIN tramite_version v ON v.tramite_id = t.id
		WHERE t.slug = $1
		  AND t.estado = 'activo'
		  AND i.estado = 'activa'
		  AND v.estado_editorial = 'publicada'
		  AND (v.valido_desde IS NULL OR v.valido_desde <= NOW())
		  AND (v.valido_hasta IS NULL OR v.valido_hasta > NOW())
		LIMIT 1
	`, slug).Scan(
		&ficha.ID, &ficha.Slug, &ficha.Titulo, &ficha.Resumen, &ficha.Descripcion,
		&ficha.PublicoObjetivo, &ficha.ResultadoEsperado, &ficha.Advertencias,
		&ficha.PlazoTexto, &ficha.VigenciaResultadoTexto,
		&ficha.RequiereCita, &ficha.URLInicio,
		&ficha.VerificadoEn, &ficha.ProximaRevisionEn,
		&inst.Codigo, &inst.Nombre, &inst.Sigla, &inst.SitioWeb,
		&ficha.Categoria,
		&versionID,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "TRAMITE_NOT_FOUND",
			"message": "Trámite no encontrado o sin versión publicada vigente",
		})
		return
	}

	ficha.Institucion = inst
	ficha.RequiereVerificacion = ficha.ProximaRevisionEn != nil && time.Now().After(*ficha.ProximaRevisionEn)

	// Modalidades
	modalidades, err := cargarModalidades(context.Background(), versionID)
	if err == nil {
		ficha.Modalidades = modalidades
	}

	// Requisitos (solo los generales, sin filtrar por modalidad aquí)
	requisitos, err := cargarRequisitos(context.Background(), versionID)
	if err == nil {
		ficha.Requisitos = requisitos
	}

	// Pasos
	pasos, err := cargarPasos(context.Background(), versionID)
	if err == nil {
		ficha.Pasos = pasos
	}

	// Costos
	costos, err := cargarCostos(context.Background(), versionID)
	if err == nil {
		ficha.Costos = costos
	}

	// Resultados
	resultados, err := cargarResultados(context.Background(), versionID)
	if err == nil {
		ficha.Resultados = resultados
	}

	// Fuentes
	fuentes, err := cargarFuentes(context.Background(), versionID)
	if err == nil {
		ficha.Fuentes = fuentes
	}

	c.JSON(http.StatusOK, ficha)
}

// GetOficinasDeModalidad godoc
// @Summary      Oficinas habilitadas para un trámite
// @Description  Devuelve las oficinas que atienden el trámite indicado, con filtro opcional por municipio.
// @Tags         Oficinas
// @Produce      json
// @Param        slug      path   string  true  "Slug del trámite"
// @Param        municipio query  int     false "ID de municipio"
// @Success      200  {object} map[string]interface{}
// @Failure      404  {object} map[string]interface{}
// @Router       /api/v1/tramites/{slug}/oficinas [get]
func GetOficinasDeModalidad(c *gin.Context) {
	slug := c.Param("slug")
	municipioStr := c.Query("municipio")

	args := []interface{}{slug}
	municipioFilter := ""
	if id, err := parseInt64(municipioStr); err == nil {
		municipioFilter = " AND o.municipio_id = $2"
		args = append(args, id)
	}

	rows, err := db.Pool.Query(context.Background(), `
		SELECT DISTINCT o.id, o.nombre, o.tipo, o.direccion,
		    o.latitud, o.longitud, o.requiere_cita, o.url_cita,
		    o.estado, m.nombre AS municipio, d.nombre AS departamento
		FROM tramite t
		JOIN tramite_version v ON v.tramite_id = t.id
		JOIN modalidad_tramite mt ON mt.tramite_version_id = v.id
		JOIN tramite_oficina tof ON tof.modalidad_id = mt.id
		JOIN oficina o ON o.id = tof.oficina_id
		LEFT JOIN municipio m ON m.id = o.municipio_id
		LEFT JOIN departamento d ON d.id = m.departamento_id
		WHERE t.slug = $1
		  AND t.estado = 'activo'
		  AND v.estado_editorial = 'publicada'
		  AND (v.valido_hasta IS NULL OR v.valido_hasta > NOW())
		  AND o.estado IN ('activa', 'temporal')
		`+municipioFilter+`
		ORDER BY o.nombre
	`, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error consultando oficinas"})
		return
	}
	defer rows.Close()

	oficinas := make([]db.OficinaSummary, 0)
	for rows.Next() {
		var o db.OficinaSummary
		if err := rows.Scan(
			&o.ID, &o.Nombre, &o.Tipo, &o.Direccion,
			&o.Latitud, &o.Longitud, &o.RequiereCita, &o.URLCita,
			&o.Estado, &o.Municipio, &o.Departamento,
		); err != nil {
			continue
		}
		oficinas = append(oficinas, o)
	}

	c.JSON(http.StatusOK, gin.H{"datos": oficinas})
}

// ─── helpers internos ──────────────────────────────────────────────────────────

func cargarModalidades(ctx context.Context, versionID int64) ([]db.ModalidadDetalle, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, tipo, nombre, descripcion, url_inicio, requiere_cita, orden
		FROM modalidad_tramite
		WHERE tramite_version_id = $1
		ORDER BY orden
	`, versionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []db.ModalidadDetalle
	for rows.Next() {
		var m db.ModalidadDetalle
		if err := rows.Scan(&m.ID, &m.Tipo, &m.Nombre, &m.Descripcion, &m.URLInicio, &m.RequiereCita, &m.Orden); err != nil {
			continue
		}

		// Cargar oficinas por modalidad
		oRows, err := db.Pool.Query(ctx, `
			SELECT o.id, o.nombre, o.tipo, o.direccion, o.latitud, o.longitud,
			       o.requiere_cita, o.url_cita, o.estado,
			       mu.nombre AS municipio, d.nombre AS departamento
			FROM tramite_oficina tof
			JOIN oficina o ON o.id = tof.oficina_id
			LEFT JOIN municipio mu ON mu.id = o.municipio_id
			LEFT JOIN departamento d ON d.id = mu.departamento_id
			WHERE tof.modalidad_id = $1
			  AND o.estado IN ('activa', 'temporal')
		`, m.ID)
		if err == nil {
			for oRows.Next() {
				var o db.OficinaSummary
				if err := oRows.Scan(&o.ID, &o.Nombre, &o.Tipo, &o.Direccion, &o.Latitud, &o.Longitud, &o.RequiereCita, &o.URLCita, &o.Estado, &o.Municipio, &o.Departamento); err == nil {
					m.Oficinas = append(m.Oficinas, o)
				}
			}
			oRows.Close()
		}

		resultado = append(resultado, m)
	}
	return resultado, nil
}

func cargarRequisitos(ctx context.Context, versionID int64) ([]db.RequisitoTramite, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, tramite_version_id, modalidad_id, codigo, titulo, descripcion,
		       obligatorio, cantidad_originales, cantidad_copias, formato,
		       vigencia_documento, emisor, grupo_alternativa, aplica_si, orden
		FROM requisito_tramite
		WHERE tramite_version_id = $1
		ORDER BY orden
	`, versionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []db.RequisitoTramite
	for rows.Next() {
		var r db.RequisitoTramite
		if err := rows.Scan(
			&r.ID, &r.TramiteVersionID, &r.ModalidadID, &r.Codigo, &r.Titulo, &r.Descripcion,
			&r.Obligatorio, &r.CantidadOriginales, &r.CantidadCopias, &r.Formato,
			&r.VigenciaDocumento, &r.Emisor, &r.GrupoAlternativa, &r.AplicaSi, &r.Orden,
		); err != nil {
			continue
		}
		resultado = append(resultado, r)
	}
	return resultado, nil
}

func cargarPasos(ctx context.Context, versionID int64) ([]db.PasoTramite, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, tramite_version_id, modalidad_id, numero, titulo, descripcion, lugar, url_accion, es_fuera_sistema
		FROM paso_tramite
		WHERE tramite_version_id = $1
		ORDER BY numero
	`, versionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []db.PasoTramite
	for rows.Next() {
		var p db.PasoTramite
		if err := rows.Scan(&p.ID, &p.TramiteVersionID, &p.ModalidadID, &p.Numero, &p.Titulo, &p.Descripcion, &p.Lugar, &p.URLAccion, &p.EsFueraSistema); err != nil {
			continue
		}
		resultado = append(resultado, p)
	}
	return resultado, nil
}

func cargarCostos(ctx context.Context, versionID int64) ([]db.CostoTramite, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, tramite_version_id, modalidad_id, concepto, moneda,
		       monto, monto_desde, monto_hasta, es_gratuito,
		       periodicidad, medio_pago, instrucciones, aplica_si, verificado_en
		FROM costo_tramite
		WHERE tramite_version_id = $1
	`, versionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []db.CostoTramite
	for rows.Next() {
		var ct db.CostoTramite
		if err := rows.Scan(
			&ct.ID, &ct.TramiteVersionID, &ct.ModalidadID, &ct.Concepto, &ct.Moneda,
			&ct.Monto, &ct.MontoDesde, &ct.MontoHasta, &ct.EsGratuito,
			&ct.Periodicidad, &ct.MedioPago, &ct.Instrucciones, &ct.AplicaSi, &ct.VerificadoEn,
		); err != nil {
			continue
		}
		resultado = append(resultado, ct)
	}
	return resultado, nil
}

func cargarResultados(ctx context.Context, versionID int64) ([]db.ResultadoTramite, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, tramite_version_id, nombre, formato, vigencia, entrega, orden
		FROM resultado_tramite
		WHERE tramite_version_id = $1
		ORDER BY orden
	`, versionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []db.ResultadoTramite
	for rows.Next() {
		var r db.ResultadoTramite
		if err := rows.Scan(&r.ID, &r.TramiteVersionID, &r.Nombre, &r.Formato, &r.Vigencia, &r.Entrega, &r.Orden); err != nil {
			continue
		}
		resultado = append(resultado, r)
	}
	return resultado, nil
}

func cargarFuentes(ctx context.Context, versionID int64) ([]db.FuenteResumen, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT fo.tipo, fo.titulo, fo.url, vf.uso
		FROM tramite_version_fuente vf
		JOIN fuente_oficial fo ON fo.id = vf.fuente_id
		WHERE vf.tramite_version_id = $1
		  AND fo.estado IN ('vigente', 'por_verificar')
		ORDER BY vf.uso
	`, versionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []db.FuenteResumen
	for rows.Next() {
		var f db.FuenteResumen
		if err := rows.Scan(&f.Tipo, &f.Titulo, &f.URL, &f.Uso); err != nil {
			continue
		}
		resultado = append(resultado, f)
	}
	return resultado, nil
}


