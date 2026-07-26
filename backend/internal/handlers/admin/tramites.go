package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"mi-tramite-bolivia-backend/internal/auditoria"
	"mi-tramite-bolivia-backend/internal/db"
	"mi-tramite-bolivia-backend/internal/fragmentador"

	"github.com/gin-gonic/gin"
)

// ─── Trámites ─────────────────────────────────────────────────────────────────

// ListarTramitesAdmin godoc
// @Summary      Listar trámites (admin)
// @Description  Lista todos los trámites con su estado editorial más reciente.
// @Tags         Admin - Trámites
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object} map[string]interface{}
// @Router       /api/v1/admin/tramites [get]
func ListarTramitesAdmin(c *gin.Context) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT t.id, t.codigo, t.slug, i.nombre AS institucion,
		       t.institucion_id, t.categoria_id, t.alcance, t.codigo_oficial, t.estado,
		       v.id AS version_id, v.numero_version,
		       v.estado_editorial, v.titulo,
		       v.verificado_en, v.proxima_revision_en,
		       v.creado_en
		FROM tramite t
		JOIN institucion i ON i.id = t.institucion_id
		LEFT JOIN LATERAL (
		    SELECT id, numero_version, estado_editorial, titulo,
		           verificado_en, proxima_revision_en, creado_en
		    FROM tramite_version
		    WHERE tramite_id = t.id
		    ORDER BY numero_version DESC
		    LIMIT 1
		) v ON true
		ORDER BY t.id DESC
		LIMIT 200
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error consultando trámites"})
		return
	}
	defer rows.Close()

	type Item struct {
		ID                int64      `json:"id"`
		Codigo            string     `json:"codigo"`
		Slug              string     `json:"slug"`
		Institucion       string     `json:"institucion"`
		InstitucionID     int64      `json:"institucion_id"`
		CategoriaID       *int64     `json:"categoria_id,omitempty"`
		Alcance           string     `json:"alcance"`
		CodigoOficial     *string    `json:"codigo_oficial,omitempty"`
		Estado            string     `json:"estado"`
		VersionID         *int64     `json:"version_id,omitempty"`
		NumeroVersion     *int       `json:"numero_version,omitempty"`
		EstadoEditorial   *string    `json:"estado_editorial,omitempty"`
		Titulo            *string    `json:"titulo,omitempty"`
		VerificadoEn      *time.Time `json:"verificado_en,omitempty"`
		ProximaRevisionEn *time.Time `json:"proxima_revision_en,omitempty"`
		CreadoEn          *time.Time `json:"creado_en,omitempty"`
	}

	datos := make([]Item, 0)
	for rows.Next() {
		var it Item
		if err := rows.Scan(
			&it.ID, &it.Codigo, &it.Slug, &it.Institucion, &it.InstitucionID, &it.CategoriaID, &it.Alcance, &it.CodigoOficial, &it.Estado,
			&it.VersionID, &it.NumeroVersion, &it.EstadoEditorial, &it.Titulo,
			&it.VerificadoEn, &it.ProximaRevisionEn, &it.CreadoEn,
		); err != nil {
			continue
		}
		datos = append(datos, it)
	}

	c.JSON(http.StatusOK, gin.H{"datos": datos})
}

// ObtenerTramiteAdmin godoc
// @Summary      Obtener trámite por ID
// @Description  Obtiene el detalle completo de un trámite y su versión más reciente.
// @Tags         Admin - Trámites
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del trámite"
// @Success      200  {object} map[string]interface{}
// @Router       /api/v1/admin/tramites/{id} [get]
func ObtenerTramiteAdmin(c *gin.Context) {
	idStr := c.Param("id")
	var id int64
	if _, err := fmt.Sscan(idStr, &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": "ID inválido"})
		return
	}

	var it struct {
		ID                int64      `json:"id"`
		Codigo            string     `json:"codigo"`
		Slug              string     `json:"slug"`
		Institucion       string     `json:"institucion"`
		Estado            string     `json:"estado"`
		VersionID         *int64     `json:"version_id,omitempty"`
		NumeroVersion     *int       `json:"numero_version,omitempty"`
		EstadoEditorial   *string    `json:"estado_editorial,omitempty"`
		Titulo            *string    `json:"titulo,omitempty"`
		Resumen           *string    `json:"resumen,omitempty"`
		Descripcion       *string    `json:"descripcion,omitempty"`
		PublicoObjetivo   *string    `json:"publico_objetivo,omitempty"`
		ResultadoEsperado *string    `json:"resultado_esperado,omitempty"`
		Advertencias      *string    `json:"advertencias,omitempty"`
		PlazoTexto        *string    `json:"plazo_texto,omitempty"`
		URLInicio         *string    `json:"url_inicio,omitempty"`
		VerificadoEn      *time.Time `json:"verificado_en,omitempty"`
		ProximaRevisionEn *time.Time `json:"proxima_revision_en,omitempty"`
		CreadoEn          *time.Time `json:"creado_en,omitempty"`
	}

	err := db.Pool.QueryRow(context.Background(), `
		SELECT t.id, t.codigo, t.slug, i.nombre AS institucion,
		       t.estado,
		       v.id AS version_id, v.numero_version,
		       v.estado_editorial, v.titulo, v.resumen, v.descripcion,
		       v.publico_objetivo, v.resultado_esperado, v.advertencias,
		       v.plazo_texto, v.url_inicio,
		       v.verificado_en, v.proxima_revision_en,
		       v.creado_en
		FROM tramite t
		JOIN institucion i ON i.id = t.institucion_id
		LEFT JOIN LATERAL (
		    SELECT id, numero_version, estado_editorial, titulo, resumen, descripcion,
		           publico_objetivo, resultado_esperado, advertencias, plazo_texto, url_inicio,
		           verificado_en, proxima_revision_en, creado_en
		    FROM tramite_version
		    WHERE tramite_id = t.id
		    ORDER BY numero_version DESC
		    LIMIT 1
		) v ON true
		WHERE t.id = $1
	`, id).Scan(
		&it.ID, &it.Codigo, &it.Slug, &it.Institucion, &it.Estado,
		&it.VersionID, &it.NumeroVersion, &it.EstadoEditorial, &it.Titulo,
		&it.Resumen, &it.Descripcion, &it.PublicoObjetivo, &it.ResultadoEsperado,
		&it.Advertencias, &it.PlazoTexto, &it.URLInicio,
		&it.VerificadoEn, &it.ProximaRevisionEn, &it.CreadoEn,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "Trámite no encontrado"})
		return
	}

	c.JSON(http.StatusOK, it)
}

// CrearTramiteAdmin godoc
// @Summary      Crear trámite
// @Description  Crea un nuevo trámite en estado activo sin versión aún.
// @Tags         Admin - Trámites
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body map[string]interface{} true "Datos del trámite"
// @Success      201  {object} map[string]interface{}
// @Router       /api/v1/admin/tramites [post]
func CrearTramiteAdmin(c *gin.Context) {
	var req struct {
		Codigo       string `json:"codigo" binding:"required"`
		Slug         string `json:"slug" binding:"required"`
		InstitucionID int64 `json:"institucion_id" binding:"required"`
		CategoriaID  *int64 `json:"categoria_id"`
		Alcance      string `json:"alcance"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}
	if req.Alcance == "" {
		req.Alcance = "nacional"
	}

	var id int64
	err := db.Pool.QueryRow(context.Background(), `
		INSERT INTO tramite (codigo, slug, institucion_id, categoria_id, alcance)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, req.Codigo, req.Slug, req.InstitucionID, req.CategoriaID, req.Alcance).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error creando trámite: " + err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	idStr := fmt.Sprintf("%d", id)
	auditoria.Registrar(c.Request.Context(), actorID, "tramite.crear", "tramite", &idStr, nil, req, nil, nil)

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Trámite creado"})
}

// ActualizarTramiteAdmin godoc
// @Summary      Actualizar trámite
// @Description  Actualiza la metadata base del trámite (slug, código, etc.).
// @Tags         Admin - Trámites
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del trámite"
// @Router       /api/v1/admin/tramites/{id} [put]
func ActualizarTramiteAdmin(c *gin.Context) {
	idStr := c.Param("id")
	var req struct {
		Codigo             *string `json:"codigo"`
		Slug               *string `json:"slug"`
		InstitucionID      *int64  `json:"institucion_id"`
		CategoriaID        *int64  `json:"categoria_id"`
		Alcance            *string `json:"alcance"`
		CodigoOficial      *string `json:"codigo_oficial"`
		Estado             *string `json:"estado"`
		ReemplazaTramiteID *int64  `json:"reemplaza_tramite_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	_, err := db.Pool.Exec(context.Background(), `
		UPDATE tramite
		SET codigo = COALESCE($2, codigo),
		    slug = COALESCE($3, slug),
		    institucion_id = COALESCE($4, institucion_id),
		    categoria_id = COALESCE($5, categoria_id),
		    alcance = COALESCE($6, alcance),
		    codigo_oficial = COALESCE($7, codigo_oficial),
		    estado = COALESCE($8, estado),
		    reemplaza_tramite_id = COALESCE($9, reemplaza_tramite_id),
		    actualizado_en = NOW()
		WHERE id = $1
	`, idStr, req.Codigo, req.Slug, req.InstitucionID, req.CategoriaID, req.Alcance, req.CodigoOficial, req.Estado, req.ReemplazaTramiteID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error actualizando trámite: " + err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	auditoria.Registrar(c.Request.Context(), actorID, "tramite.actualizar", "tramite", &idStr, nil, req, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Trámite actualizado"})
}

// EliminarTramiteAdmin godoc
// @Summary      Eliminar trámite
// @Description  Cambia el estado del trámite a 'retirado' (borrado lógico).
// @Tags         Admin - Trámites
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del trámite"
// @Router       /api/v1/admin/tramites/{id} [delete]
func EliminarTramiteAdmin(c *gin.Context) {
	idStr := c.Param("id")
	
	_, err := db.Pool.Exec(context.Background(), `
		UPDATE tramite SET estado = 'retirado', actualizado_en = NOW() WHERE id = $1
	`, idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error eliminando trámite: " + err.Error()})
		return
	}

	actorID := actorIDFromCtx(c)
	auditoria.Registrar(c.Request.Context(), actorID, "tramite.eliminar", "tramite", &idStr, nil, nil, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Trámite retirado (inactivado)"})
}

// CrearVersion godoc
// @Summary      Crear versión borrador
// @Description  Crea una nueva versión borrador de un trámite existente.
// @Tags         Admin - Versiones
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del trámite"
// @Success      201  {object} map[string]interface{}
// @Router       /api/v1/admin/tramites/{id}/versiones [post]
func CrearVersion(c *gin.Context) {
	tramiteIDStr := c.Param("id")
	var tramiteID int64
	if _, err := fmt.Sscan(tramiteIDStr, &tramiteID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": "ID inválido"})
		return
	}

	var req struct {
		Titulo      string `json:"titulo" binding:"required"`
		Resumen     string `json:"resumen" binding:"required"`
		Descripcion string `json:"descripcion" binding:"required"`
		PublicoObjetivo  *string `json:"publico_objetivo"`
		ResultadoEsperado *string `json:"resultado_esperado"`
		Advertencias     *string `json:"advertencias"`
		PlazoTexto       *string `json:"plazo_texto"`
		URLInicio        *string `json:"url_inicio"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	// Determinar número de versión siguiente
	var ultimoNum int
	_ = db.Pool.QueryRow(context.Background(), `
		SELECT COALESCE(MAX(numero_version), 0) FROM tramite_version WHERE tramite_id = $1
	`, tramiteID).Scan(&ultimoNum)

	actorID := actorIDFromCtx(c)
	var versionID int64
	err := db.Pool.QueryRow(context.Background(), `
		INSERT INTO tramite_version
		    (tramite_id, numero_version, estado_editorial, titulo, resumen, descripcion,
		     publico_objetivo, resultado_esperado, advertencias, plazo_texto, url_inicio,
		     creado_por)
		VALUES ($1, $2, 'borrador', $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`, tramiteID, ultimoNum+1, req.Titulo, req.Resumen, req.Descripcion,
		req.PublicoObjetivo, req.ResultadoEsperado, req.Advertencias,
		req.PlazoTexto, req.URLInicio, actorID,
	).Scan(&versionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error creando versión: " + err.Error()})
		return
	}

	idStr := fmt.Sprintf("%d", versionID)
	auditoria.Registrar(c.Request.Context(), actorID, "version.crear", "tramite_version", &idStr, nil, req, nil, nil)

	c.JSON(http.StatusCreated, gin.H{"version_id": versionID, "numero_version": ultimoNum + 1})
}

// ActualizarVersion godoc
// @Summary      Editar versión borrador
// @Description  Actualiza los campos de una versión en estado borrador. No permite editar versiones publicadas.
// @Tags         Admin - Versiones
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID de la versión"
// @Success      200  {object} map[string]interface{}
// @Router       /api/v1/admin/versiones/{id} [put]
func ActualizarVersion(c *gin.Context) {
	versionIDStr := c.Param("id")
	var versionID int64
	if _, err := fmt.Sscan(versionIDStr, &versionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": "ID inválido"})
		return
	}

	// Verificar que está en estado editable
	var estado string
	if err := db.Pool.QueryRow(context.Background(), `
		SELECT estado_editorial FROM tramite_version WHERE id = $1
	`, versionID).Scan(&estado); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "Versión no encontrada"})
		return
	}
	if estado == "publicada" || estado == "reemplazada" {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    "IMMUTABLE_VERSION",
			"message": "No se puede editar una versión publicada. Cree una nueva versión.",
		})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	// Actualizar solo campos seguros (whitelist)
	campos := map[string]interface{}{}
	allowed := []string{"titulo", "resumen", "descripcion", "publico_objetivo",
		"resultado_esperado", "advertencias", "plazo_texto", "url_inicio",
		"requiere_cita", "notas_internas", "vigencia_resultado_texto"}
	for _, k := range allowed {
		if v, ok := req[k]; ok {
			campos[k] = v
		}
	}

	if len(campos) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": "No hay campos válidos para actualizar"})
		return
	}

	// Construir SET dinámico
	setSQL := ""
	args := []interface{}{}
	idx := 1
	for k, v := range campos {
		if setSQL != "" {
			setSQL += ", "
		}
		setSQL += fmt.Sprintf("%s = $%d", k, idx)
		args = append(args, v)
		idx++
	}
	args = append(args, versionID)

	_, err := db.Pool.Exec(context.Background(),
		"UPDATE tramite_version SET "+setSQL+" WHERE id = $"+fmt.Sprintf("%d", idx),
		args...,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error actualizando versión"})
		return
	}

	actorID := actorIDFromCtx(c)
	idStr := versionIDStr
	auditoria.Registrar(c.Request.Context(), actorID, "version.actualizar", "tramite_version", &idStr, nil, campos, nil, nil)

	c.JSON(http.StatusOK, gin.H{"message": "Versión actualizada"})
}

// EnviarARevision godoc
// @Summary      Enviar versión a revisión
// @Tags         Admin - Versiones
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID de la versión"
// @Success      200  {object} map[string]interface{}
// @Router       /api/v1/admin/versiones/{id}/enviar-revision [post]
func EnviarARevision(c *gin.Context) {
	cambiarEstadoVersion(c, "borrador", "en_revision", "version.enviar_revision")
}

// AprobarVersion godoc
// @Summary      Aprobar versión
// @Tags         Admin - Versiones
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID de la versión"
// @Success      200  {object} map[string]interface{}
// @Router       /api/v1/admin/versiones/{id}/aprobar [post]
func AprobarVersion(c *gin.Context) {
	cambiarEstadoVersion(c, "en_revision", "en_revision", "version.aprobar") // actualiza revisado_por
}

// RechazarVersion godoc
// @Summary      Rechazar versión
// @Tags         Admin - Versiones
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID de la versión"
// @Success      200  {object} map[string]interface{}
// @Router       /api/v1/admin/versiones/{id}/rechazar [post]
func RechazarVersion(c *gin.Context) {
	versionIDStr := c.Param("id")
	var versionID int64
	if _, err := fmt.Sscan(versionIDStr, &versionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": "ID inválido"})
		return
	}

	var req struct {
		Observaciones string `json:"observaciones"`
	}
	_ = c.ShouldBindJSON(&req)

	actorID := actorIDFromCtx(c)
	_, err := db.Pool.Exec(context.Background(), `
		UPDATE tramite_version
		SET estado_editorial = 'rechazada',
		    revisado_por = $2,
		    notas_internas = COALESCE(notas_internas, '') || E'\n[RECHAZADA]: ' || $3
		WHERE id = $1 AND estado_editorial = 'en_revision'
	`, versionID, actorID, req.Observaciones)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error rechazando versión"})
		return
	}

	idStr := versionIDStr
	auditoria.Registrar(c.Request.Context(), actorID, "version.rechazar", "tramite_version", &idStr, nil, req, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Versión rechazada"})
}

// PublicarVersion godoc
// @Summary      Publicar versión
// @Description  Publica la versión aprobada en transacción: cierra la anterior, publica la nueva y genera fragmentos de conocimiento.
// @Tags         Admin - Versiones
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID de la versión"
// @Success      200  {object} map[string]interface{}
// @Router       /api/v1/admin/versiones/{id}/publicar [post]
func PublicarVersion(c *gin.Context) {
	versionIDStr := c.Param("id")
	var versionID int64
	if _, err := fmt.Sscan(versionIDStr, &versionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": "ID inválido"})
		return
	}

	var req struct {
		VerificadoEn      *time.Time `json:"verificado_en"`
		ProximaRevisionEn *time.Time `json:"proxima_revision_en"`
	}
	_ = c.ShouldBindJSON(&req)

	if req.VerificadoEn == nil {
		now := time.Now()
		req.VerificadoEn = &now
	}

	// Leer tramite_id y validar estado
	var tramiteID int64
	var estadoEditorial string
	var fuentes int
	err := db.Pool.QueryRow(context.Background(), `
		SELECT v.tramite_id, v.estado_editorial,
		       COUNT(vf.fuente_id) AS fuentes
		FROM tramite_version v
		LEFT JOIN tramite_version_fuente vf ON vf.tramite_version_id = v.id
		WHERE v.id = $1
		GROUP BY v.tramite_id, v.estado_editorial
	`, versionID).Scan(&tramiteID, &estadoEditorial, &fuentes)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "Versión no encontrada"})
		return
	}
	if estadoEditorial != "en_revision" {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    "INVALID_STATE",
			"message": "Solo se puede publicar una versión en estado en_revision",
		})
		return
	}
	if fuentes == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    "NO_FUENTE",
			"message": "La versión debe tener al menos una fuente oficial",
		})
		return
	}

	actorID := actorIDFromCtx(c)

	// Transacción: cerrar anterior + publicar nueva
	tx, err := db.Pool.Begin(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error iniciando transacción"})
		return
	}
	defer tx.Rollback(context.Background())

	// Marcar versión anterior como reemplazada
	_, err = tx.Exec(context.Background(), `
		UPDATE tramite_version
		SET estado_editorial = 'reemplazada',
		    valido_hasta = NOW()
		WHERE tramite_id = $1
		  AND estado_editorial = 'publicada'
		  AND valido_hasta IS NULL
	`, tramiteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error cerrando versión anterior"})
		return
	}

	// Publicar nueva versión
	_, err = tx.Exec(context.Background(), `
		UPDATE tramite_version
		SET estado_editorial = 'publicada',
		    publicado_en = NOW(),
		    valido_desde = NOW(),
		    verificado_en = $2,
		    proxima_revision_en = $3,
		    aprobado_por = $4
		WHERE id = $1
	`, versionID, req.VerificadoEn, req.ProximaRevisionEn, actorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error publicando versión: " + err.Error()})
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error confirmando transacción"})
		return
	}

	// Generar fragmentos de conocimiento en background
	go func() {
		ctx := context.Background()
		if err := fragmentador.Fragmentar(ctx, versionID, "text-embedding-004"); err != nil {
			// Log pero no fallar la respuesta
			fmt.Printf("[publicar] error fragmentando versión %d: %v\n", versionID, err)
		}
	}()

	idStr := versionIDStr
	auditoria.Registrar(c.Request.Context(), actorID, "version.publicar", "tramite_version", &idStr, nil, nil, nil, nil)

	c.JSON(http.StatusOK, gin.H{
		"message":    "Versión publicada exitosamente",
		"version_id": versionID,
	})
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func cambiarEstadoVersion(c *gin.Context, estadoEsperado, nuevoEstado, accion string) {
	versionIDStr := c.Param("id")
	var versionID int64
	if _, err := fmt.Sscan(versionIDStr, &versionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": "ID inválido"})
		return
	}

	actorID := actorIDFromCtx(c)
	var result int64
	err := db.Pool.QueryRow(context.Background(), `
		UPDATE tramite_version
		SET estado_editorial = $3, revisado_por = $4
		WHERE id = $1 AND estado_editorial = $2
		RETURNING id
	`, versionID, estadoEsperado, nuevoEstado, actorID).Scan(&result)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "Versión no encontrada o en estado incorrecto"})
		return
	}

	idStr := versionIDStr
	auditoria.Registrar(c.Request.Context(), actorID, accion, "tramite_version", &idStr, nil, nil, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Estado actualizado"})
}

func actorIDFromCtx(c *gin.Context) *int64 {
	v, ok := c.Get("usuario_id")
	if !ok {
		return nil
	}
	id, ok := v.(int64)
	if !ok {
		return nil
	}
	return &id
}

func marshalJSON(v interface{}) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
