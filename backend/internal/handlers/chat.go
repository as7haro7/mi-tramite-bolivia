package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"strings"
	"time"

	"mi-tramite-bolivia-backend/internal/ai"
	"mi-tramite-bolivia-backend/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/pgvector/pgvector-go"
)

type CrearConversacionRequest struct {
	MunicipioID *int64 `json:"municipio_id"`
	Idioma      string `json:"idioma"`
}

type EnviarMensajeRequest struct {
	Contenido   string `json:"contenido" binding:"required,max=2000"`
	MunicipioID *int64 `json:"municipio_id"`
}

type FeedbackRequest struct {
	Valoracion int    `json:"valoracion" binding:"required"`  // 1 o -1
	Motivo     string `json:"motivo"`
	Comentario string `json:"comentario"`
}

const umbralSimilitudRAG = 0.70

// CrearConversacion godoc
// @Summary      Iniciar conversación
// @Description  Crea una sesión anónima y una conversación de chat.
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        request body CrearConversacionRequest false "Contexto inicial"
// @Success      201  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/chat/conversaciones [post]
func CrearConversacion(c *gin.Context) {
	var req CrearConversacionRequest
	_ = c.ShouldBindJSON(&req)

	if req.Idioma == "" {
		req.Idioma = "es-BO"
	}

	// Hashear IP para privacidad
	ipRaw := c.ClientIP()
	idHash := hashear(ipRaw + time.Now().String())

	expira := time.Now().Add(30 * 24 * time.Hour)

	var sesionID string
	err := db.Pool.QueryRow(context.Background(), `
		INSERT INTO sesion_anonima
		    (identificador_hash, idioma, municipio_id, expira_en)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, idHash, req.Idioma, req.MunicipioID, expira).Scan(&sesionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "SESSION_ERROR",
			"message": "Error creando sesión",
		})
		return
	}

	var convID string
	err = db.Pool.QueryRow(context.Background(), `
		INSERT INTO conversacion (sesion_id) VALUES ($1) RETURNING id
	`, sesionID).Scan(&convID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "CONVERSATION_ERROR",
			"message": "Error creando conversación",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"conversacion_id": convID,
		"sesion_id":       sesionID,
		"expira_en":       expira,
	})
}

// EnviarMensaje godoc
// @Summary      Enviar mensaje en conversación
// @Description  Procesa la consulta usando RAG sobre fragmentos publicados y responde con citas.
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        id   path string true "ID de la conversación"
// @Param        request body EnviarMensajeRequest true "Mensaje del usuario"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      404  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/chat/conversaciones/{id}/mensajes [post]
func EnviarMensaje(c *gin.Context) {
	convID := c.Param("id")
	startTime := time.Now()

	var req EnviarMensajeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "VALIDATION_ERROR",
			"message": "Mensaje inválido",
			"details": err.Error(),
		})
		return
	}

	// Sanitizar contra prompt injection básico
	if strings.Contains(strings.ToLower(req.Contenido), "ignore previous") ||
		strings.Contains(strings.ToLower(req.Contenido), "system:") {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "SECURITY_BLOCK",
			"message": "Mensaje no permitido",
		})
		return
	}

	// Verificar que la conversación existe y está activa
	var sesionID string
	var convEstado string
	err := db.Pool.QueryRow(context.Background(), `
		SELECT c.sesion_id, c.estado FROM conversacion c
		JOIN sesion_anonima s ON s.id = c.sesion_id
		WHERE c.id = $1 AND c.estado = 'activa' AND s.expira_en > NOW()
	`, convID).Scan(&sesionID, &convEstado)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "CONVERSATION_NOT_FOUND",
			"message": "Conversación no encontrada o expirada",
		})
		return
	}

	// Guardar mensaje del usuario
	var msgUsuarioID string
	err = db.Pool.QueryRow(context.Background(), `
		INSERT INTO mensaje_conversacion (conversacion_id, rol, contenido)
		VALUES ($1, 'usuario', $2) RETURNING id
	`, convID, req.Contenido).Scan(&msgUsuarioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error guardando mensaje"})
		return
	}

	// Actualizar última actividad
	_, _ = db.Pool.Exec(context.Background(), `
		UPDATE sesion_anonima SET ultima_actividad_en = NOW() WHERE id = $1
	`, sesionID)

	// Generar embedding de la consulta
	vector, embErr := ai.GenerateEmbedding(req.Contenido)
	var respuesta string
	var citas []db.CitaRAG

	if embErr != nil {
		// Degradar a búsqueda por texto si el LLM falla (RF-CHAT-010)
		respuesta = fallbackBusqueda(context.Background(), req.Contenido)
	} else {
		queryVec := pgvector.NewVector(vector)

		// Recuperar fragmentos de versiones publicadas vigentes con umbral
		fragRows, err := db.Pool.Query(context.Background(), `
			SELECT fk.id, fk.contenido, fk.tipo,
			       t.slug, v.titulo,
			       (fk.embedding <=> $1) AS distancia
			FROM fragmento_conocimiento fk
			JOIN tramite_version v ON v.id = fk.tramite_version_id
			JOIN tramite t ON t.id = v.tramite_id
			WHERE fk.estado_embedding = 'listo'
			  AND v.estado_editorial = 'publicada'
			  AND (v.valido_hasta IS NULL OR v.valido_hasta > NOW())
			  AND t.estado = 'activo'
			ORDER BY fk.embedding <=> $1
			LIMIT 5
		`, queryVec)

		var fragmentos []struct {
			id        string
			contenido string
			tipo      string
			slug      string
			titulo    string
			distancia float64
		}

		if err == nil {
			for fragRows.Next() {
				var f struct {
					id        string
					contenido string
					tipo      string
					slug      string
					titulo    string
					distancia float64
				}
				if err := fragRows.Scan(&f.id, &f.contenido, &f.tipo, &f.slug, &f.titulo, &f.distancia); err == nil {
					fragmentos = append(fragmentos, f)
				}
			}
			fragRows.Close()
		}

		// Guardar recuperacion_rag
		var recuperacionID string
		_ = db.Pool.QueryRow(context.Background(), `
			INSERT INTO recuperacion_rag (mensaje_usuario_id, consulta_normalizada, modelo_embedding)
			VALUES ($1, $2, 'text-embedding-004')
			RETURNING id
		`, msgUsuarioID, req.Contenido).Scan(&recuperacionID)

		// Filtrar por umbral y construir contexto
		var contextoPartes []string
		for pos, f := range fragmentos {
			similitud := 1 - f.distancia
			incluidoEnPrompt := similitud >= umbralSimilitudRAG

			if recuperacionID != "" {
				_, _ = db.Pool.Exec(context.Background(), `
					INSERT INTO cita_rag (recuperacion_id, fragmento_id, posicion, similitud, incluida_en_prompt)
					VALUES ($1, $2, $3, $4, $5)
					ON CONFLICT DO NOTHING
				`, recuperacionID, f.id, pos+1, similitud, incluidoEnPrompt)
			}

			if incluidoEnPrompt {
				contextoPartes = append(contextoPartes, f.contenido)
				citas = append(citas, db.CitaRAG{
					FragmentoID:      f.id,
					Posicion:         pos + 1,
					Similitud:        &similitud,
					IncluidaEnPrompt: incluidoEnPrompt,
					TramiteSlug:      &f.slug,
					TramiteTitulo:    &f.titulo,
					TipoFragmento:    &f.tipo,
				})
			}
		}

		if len(contextoPartes) == 0 {
			respuesta = "No tengo información oficial verificada suficiente para responder eso. Puedo ayudarte a buscar el trámite por institución o categoría."
		} else {
			contexto := strings.Join(contextoPartes, "\n\n---\n\n")
			respuesta, _ = ai.GenerateChatResponse(contexto, req.Contenido)
			if respuesta == "" {
				respuesta = fallbackBusqueda(context.Background(), req.Contenido)
			}
		}
	}

	latenciaMs := int(time.Since(startTime).Milliseconds())

	// Guardar respuesta del asistente
	modelo := "gemini-2.0-flash"
	var msgAsistID string
	_ = db.Pool.QueryRow(context.Background(), `
		INSERT INTO mensaje_conversacion
		    (conversacion_id, rol, contenido, modelo, latencia_ms)
		VALUES ($1, 'asistente', $2, $3, $4)
		RETURNING id
	`, convID, respuesta, modelo, latenciaMs).Scan(&msgAsistID)

	if citas == nil {
		citas = []db.CitaRAG{}
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje_id":  msgAsistID,
		"respuesta":   respuesta,
		"citas":       citas,
		"latencia_ms": latenciaMs,
	})
}

// RegistrarFeedback godoc
// @Summary      Valorar respuesta del asistente
// @Description  Registra si la respuesta fue útil o no.
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        id   path string true "ID del mensaje del asistente"
// @Param        request body FeedbackRequest true "Valoración"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Router       /api/v1/mensajes/{id}/feedback [post]
func RegistrarFeedback(c *gin.Context) {
	msgID := c.Param("id")
	var req FeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	if req.Valoracion != 1 && req.Valoracion != -1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": "Valoración debe ser 1 o -1"})
		return
	}

	_, err := db.Pool.Exec(context.Background(), `
		INSERT INTO feedback_respuesta
		    (mensaje_asistente_id, valoracion, motivo, comentario)
		VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''))
		ON CONFLICT (mensaje_asistente_id) DO UPDATE
		    SET valoracion = EXCLUDED.valoracion,
		        motivo = EXCLUDED.motivo,
		        comentario = EXCLUDED.comentario
	`, msgID, req.Valoracion, req.Motivo, req.Comentario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error guardando feedback"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feedback registrado"})
}

// fallbackBusqueda devuelve un mensaje genérico con búsqueda textual cuando falla el LLM.
func fallbackBusqueda(ctx context.Context, consulta string) string {
	normalizado := normalizarTexto(consulta)
	var titulo, slug string
	err := db.Pool.QueryRow(ctx, `
		SELECT v.titulo, t.slug
		FROM tramite t
		JOIN tramite_version v ON v.tramite_id = t.id
		WHERE t.estado = 'activo'
		  AND v.estado_editorial = 'publicada'
		  AND public.normalizar_busqueda(v.titulo) ILIKE '%'||$1||'%'
		LIMIT 1
	`, normalizado).Scan(&titulo, &slug)
	if err == nil {
		return "Encontré información sobre \"" + titulo + "\". Puedes consultar la ficha detallada de este trámite desde el buscador de la aplicación."
	}
	return "El asistente conversacional se encuentra temporalmente fuera de línea. Puedes buscar tu trámite directamente en el catálogo oficial o consultar a la institución correspondiente."
}

func hashear(s string) string {
	secret := os.Getenv("IP_HASH_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(s))
	return hex.EncodeToString(mac.Sum(nil))
}
