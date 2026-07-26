package handlers

import (
	"context"
	"net/http"
	"strconv"

	"mi-tramite-bolivia-backend/internal/db"

	"github.com/gin-gonic/gin"
)

// GetInstituciones godoc
// @Summary      Listar instituciones activas
// @Description  Devuelve todas las instituciones activas con sus alias.
// @Tags         Catálogo
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/instituciones [get]
func GetInstituciones(c *gin.Context) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT i.id, i.codigo, i.nombre, i.sigla, i.tipo, i.sitio_web, i.estado,
		       am.url_publica AS logo_url
		FROM institucion i
		LEFT JOIN archivo_multimedia am ON am.id = i.logo_archivo_id
		    AND am.variante = 'original' AND am.estado = 'activo'
		WHERE i.estado = 'activa'
		ORDER BY i.nombre
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error consultando instituciones"})
		return
	}
	defer rows.Close()

	type InstItem struct {
		ID       int64   `json:"id"`
		Codigo   string  `json:"codigo"`
		Nombre   string  `json:"nombre"`
		Sigla    string  `json:"sigla"`
		Tipo     string  `json:"tipo"`
		SitioWeb *string `json:"sitio_web,omitempty"`
		Estado   string  `json:"estado"`
		LogoURL  *string `json:"logo_url,omitempty"`
	}

	datos := make([]InstItem, 0)
	for rows.Next() {
		var it InstItem
		if err := rows.Scan(&it.ID, &it.Codigo, &it.Nombre, &it.Sigla, &it.Tipo, &it.SitioWeb, &it.Estado, &it.LogoURL); err != nil {
			continue
		}
		datos = append(datos, it)
	}

	c.JSON(http.StatusOK, gin.H{"datos": datos})
}

// GetCategorias godoc
// @Summary      Listar categorías
// @Description  Devuelve el árbol de categorías activas ordenadas.
// @Tags         Catálogo
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/categorias [get]
func GetCategorias(c *gin.Context) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT id, padre_id, codigo, nombre, icono, orden
		FROM categoria
		WHERE activa = true
		ORDER BY orden, nombre
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error consultando categorías"})
		return
	}
	defer rows.Close()

	datos := make([]db.Categoria, 0)
	for rows.Next() {
		var cat db.Categoria
		if err := rows.Scan(&cat.ID, &cat.PadreID, &cat.Codigo, &cat.Nombre, &cat.Icono, &cat.Orden); err != nil {
			continue
		}
		cat.Activa = true
		datos = append(datos, cat)
	}

	c.JSON(http.StatusOK, gin.H{"datos": datos})
}

// GetOficinas godoc
// @Summary      Listar oficinas habilitadas
// @Description  Devuelve oficinas activas o temporales con filtro opcional por municipio.
// @Tags         Oficinas
// @Produce      json
// @Param        municipio query int false "ID de municipio"
// @Param        tipo      query string false "Tipo de oficina (oficina, ventanilla, etc.)"
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/oficinas [get]
func GetOficinas(c *gin.Context) {
	args := []interface{}{}
	idx := 1
	filtros := []string{"o.estado IN ('activa', 'temporal')"}

	if munStr := c.Query("municipio"); munStr != "" {
		if id, err := strconv.ParseInt(munStr, 10, 64); err == nil {
			filtros = append(filtros, "o.municipio_id = $"+strconv.Itoa(idx))
			args = append(args, id)
			idx++
		}
	}
	if tipo := c.Query("tipo"); tipo != "" {
		filtros = append(filtros, "o.tipo = $"+strconv.Itoa(idx))
		args = append(args, tipo)
		idx++
	}

	whereSQL := ""
	if len(filtros) > 0 {
		whereSQL = "WHERE "
		for i, f := range filtros {
			if i > 0 {
				whereSQL += " AND "
			}
			whereSQL += f
		}
	}

	query := `
		SELECT o.id, o.nombre, o.tipo, o.direccion, o.latitud, o.longitud,
		       o.requiere_cita, o.url_cita, o.estado,
		       m.nombre AS municipio, d.nombre AS departamento
		FROM oficina o
		LEFT JOIN municipio m ON m.id = o.municipio_id
		LEFT JOIN departamento d ON d.id = m.departamento_id
		JOIN institucion i ON i.id = o.institucion_id AND i.estado = 'activa'
		` + whereSQL + `
		ORDER BY o.nombre
		LIMIT 200`

	rows, err := db.Pool.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error consultando oficinas"})
		return
	}
	defer rows.Close()

	datos := make([]db.OficinaSummary, 0)
	for rows.Next() {
		var o db.OficinaSummary
		if err := rows.Scan(&o.ID, &o.Nombre, &o.Tipo, &o.Direccion, &o.Latitud, &o.Longitud, &o.RequiereCita, &o.URLCita, &o.Estado, &o.Municipio, &o.Departamento); err != nil {
			continue
		}
		datos = append(datos, o)
	}

	c.JSON(http.StatusOK, gin.H{"datos": datos})
}
