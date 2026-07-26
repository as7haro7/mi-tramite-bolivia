package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mi-tramite-bolivia-backend/internal/db"

	"github.com/gin-gonic/gin"
)

// GetTramites godoc
// @Summary      Listar trámites publicados
// @Description  Devuelve trámites con versión publicada y vigente. Soporta filtros y paginación con cursor.
// @Tags         Catálogo
// @Produce      json
// @Param        q          query  string  false  "Búsqueda por texto (nombre o alias)"
// @Param        institucion query  string  false  "Código de institución"
// @Param        categoria   query  string  false  "Código de categoría"
// @Param        municipio   query  int     false  "ID de municipio"
// @Param        modalidad   query  string  false  "Tipo de modalidad (presencial, en_linea, etc.)"
// @Param        cursor      query  int     false  "ID del último elemento (paginación)"
// @Param        limite      query  int     false  "Elementos por página (máx 50, def 20)"
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/tramites [get]
func GetTramites(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	institucion := strings.TrimSpace(c.Query("institucion"))
	categoria := strings.TrimSpace(c.Query("categoria"))
	municipioStr := c.Query("municipio")
	modalidad := strings.TrimSpace(c.Query("modalidad"))
	cursorStr := c.Query("cursor")
	limiteStr := c.Query("limite")

	limite := 20
	if n, err := strconv.Atoi(limiteStr); err == nil && n > 0 && n <= 50 {
		limite = n
	}
	var cursorID *int64
	if id, err := strconv.ParseInt(cursorStr, 10, 64); err == nil && id > 0 {
		cursorID = &id
	}
	var municipioID *int64
	if id, err := strconv.ParseInt(municipioStr, 10, 64); err == nil && id > 0 {
		municipioID = &id
	}

	args := []interface{}{}
	idx := 1

	where := []string{
		"t.estado = 'activo'",
		"i.estado = 'activa'",
		"v.estado_editorial = 'publicada'",
		"(v.valido_desde IS NULL OR v.valido_desde <= NOW())",
		"(v.valido_hasta IS NULL OR v.valido_hasta > NOW())",
	}

	if cursorID != nil {
		where = append(where, "t.id < $"+strconv.Itoa(idx))
		args = append(args, *cursorID)
		idx++
	}
	if q != "" {
		normalizado := normalizarTexto(q)
		where = append(where, `(
			public.normalizar_busqueda(v.titulo) ILIKE '%'||$`+strconv.Itoa(idx)+`||'%'
			OR public.normalizar_busqueda(v.resumen) ILIKE '%'||$`+strconv.Itoa(idx)+`||'%'
			OR EXISTS (
				SELECT 1 FROM tramite_alias ta
				WHERE ta.tramite_id = t.id
				  AND ta.normalizado ILIKE '%'||$`+strconv.Itoa(idx)+`||'%'
			)
		)`)
		args = append(args, normalizado)
		idx++
	}
	if institucion != "" {
		where = append(where, "i.codigo = $"+strconv.Itoa(idx))
		args = append(args, strings.ToUpper(institucion))
		idx++
	}
	if categoria != "" {
		where = append(where, "c.codigo = $"+strconv.Itoa(idx))
		args = append(args, categoria)
		idx++
	}
	if municipioID != nil {
		where = append(where, `EXISTS (
			SELECT 1 FROM modalidad_tramite mt
			JOIN tramite_oficina tof ON tof.modalidad_id = mt.id
			JOIN oficina o ON o.id = tof.oficina_id
			WHERE mt.tramite_version_id = v.id AND o.municipio_id = $`+strconv.Itoa(idx)+`
		)`)
		args = append(args, *municipioID)
		idx++
	}
	if modalidad != "" {
		where = append(where, `EXISTS (
			SELECT 1 FROM modalidad_tramite mt
			WHERE mt.tramite_version_id = v.id AND mt.tipo = $`+strconv.Itoa(idx)+`
		)`)
		args = append(args, modalidad)
		idx++
	}

	args = append(args, limite+1)
	whereSQL := "WHERE " + strings.Join(where, " AND ")

	query := `
		SELECT t.id, t.slug, v.titulo, v.resumen,
		       i.nombre AS institucion, i.sigla,
		       c.nombre AS categoria,
		       v.verificado_en, v.proxima_revision_en
		FROM tramite t
		JOIN institucion i ON i.id = t.institucion_id
		LEFT JOIN categoria c ON c.id = t.categoria_id
		JOIN tramite_version v ON v.tramite_id = t.id
		` + whereSQL + `
		ORDER BY t.id DESC
		LIMIT $` + strconv.Itoa(idx)

	rows, err := db.Pool.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "DB_ERROR",
			"message": "Error consultando trámites",
		})
		return
	}
	defer rows.Close()

	ahora := time.Now()
	items := make([]db.TramiteListItem, 0, limite)
	for rows.Next() {
		var item db.TramiteListItem
		if err := rows.Scan(
			&item.ID, &item.Slug, &item.Titulo, &item.Resumen,
			&item.Institucion, &item.InstitucionSigla,
			&item.Categoria,
			&item.VerificadoEn, &item.ProximaRevisionEn,
		); err != nil {
			continue
		}
		item.RequiereVerificacion = item.ProximaRevisionEn != nil && ahora.After(*item.ProximaRevisionEn)
		items = append(items, item)
	}

	hayMas := len(items) > limite
	if hayMas {
		items = items[:limite]
	}

	var nextCursor *int64
	if hayMas && len(items) > 0 {
		last := items[len(items)-1].ID
		nextCursor = &last
	}

	c.JSON(http.StatusOK, gin.H{
		"datos":       items,
		"hay_mas":     hayMas,
		"next_cursor": nextCursor,
		"limite":      limite,
	})
}

func normalizarTexto(s string) string {
	s = strings.ToLower(s)
	replacer := strings.NewReplacer(
		"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u",
		"ü", "u", "ñ", "n",
	)
	return replacer.Replace(s)
}
