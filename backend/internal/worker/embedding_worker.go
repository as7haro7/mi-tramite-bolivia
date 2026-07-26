package worker

import (
	"context"
	"log"
	"time"

	"mi-tramite-bolivia-backend/internal/ai"
	"mi-tramite-bolivia-backend/internal/db"

	"github.com/pgvector/pgvector-go"
)

const (
	maxIntentos      = 5
	tickerInterval   = 30 * time.Second
	backoffBase      = 5 * time.Second
)

// StartEmbeddingWorker lanza el bucle de procesamiento de trabajos de embedding.
// Se detiene cuando se cancela el contexto.
func StartEmbeddingWorker(ctx context.Context) {
	log.Println("[embedding-worker] iniciando")
	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[embedding-worker] detenido")
			return
		case <-ticker.C:
			procesarPendientes(ctx)
		}
	}
}

func procesarPendientes(ctx context.Context) {
	// Reclamar hasta 10 trabajos pendientes marcándolos como 'procesando'
	rows, err := db.Pool.Query(ctx, `
		UPDATE trabajo_embedding
		SET estado = 'procesando', iniciado_en = NOW()
		WHERE id IN (
		    SELECT id FROM trabajo_embedding
		    WHERE estado = 'pendiente'
		      AND disponible_desde <= NOW()
		      AND intentos < $1
		    ORDER BY disponible_desde
		    LIMIT 10
		    FOR UPDATE SKIP LOCKED
		)
		RETURNING id, fragmento_id, modelo, intentos
	`, maxIntentos)
	if err != nil {
		log.Printf("[embedding-worker] error obteniendo trabajos: %v", err)
		return
	}
	defer rows.Close()

	type trabajo struct {
		id          string
		fragmentoID string
		modelo      string
		intentos    int
	}
	var trabajos []trabajo
	for rows.Next() {
		var t trabajo
		if err := rows.Scan(&t.id, &t.fragmentoID, &t.modelo, &t.intentos); err != nil {
			continue
		}
		trabajos = append(trabajos, t)
	}
	rows.Close()

	for _, t := range trabajos {
		procesarTrabajo(ctx, t.id, t.fragmentoID, t.modelo, t.intentos)
	}
}

func procesarTrabajo(ctx context.Context, trabajoID, fragmentoID, modelo string, intentos int) {
	// Leer contenido del fragmento
	var contenido string
	err := db.Pool.QueryRow(ctx,
		`SELECT contenido FROM fragmento_conocimiento WHERE id = $1`, fragmentoID,
	).Scan(&contenido)
	if err != nil {
		marcarFallido(ctx, trabajoID, fragmentoID, intentos, "fragmento no encontrado: "+err.Error())
		return
	}

	// Generar embedding
	vector, err := ai.GenerateEmbedding(contenido)
	if err != nil {
		backoff := backoffBase * time.Duration(1<<uint(intentos))
		disponibleDesde := time.Now().Add(backoff)
		_, _ = db.Pool.Exec(ctx, `
			UPDATE trabajo_embedding
			SET estado = 'pendiente',
			    intentos = intentos + 1,
			    error_detalle = $2,
			    disponible_desde = $3
			WHERE id = $1
		`, trabajoID, err.Error(), disponibleDesde)

		if intentos+1 >= maxIntentos {
			marcarFallido(ctx, trabajoID, fragmentoID, intentos, err.Error())
		}
		return
	}

	vec := pgvector.NewVector(vector)

	// Actualizar fragmento con el vector
	_, err = db.Pool.Exec(ctx, `
		UPDATE fragmento_conocimiento
		SET embedding = $2,
		    modelo_embedding = $3,
		    estado_embedding = 'listo',
		    generado_en = NOW()
		WHERE id = $1
	`, fragmentoID, vec, modelo)
	if err != nil {
		marcarFallido(ctx, trabajoID, fragmentoID, intentos, "error guardando vector: "+err.Error())
		return
	}

	// Marcar trabajo como completado
	_, _ = db.Pool.Exec(ctx, `
		UPDATE trabajo_embedding
		SET estado = 'completado', completado_en = NOW()
		WHERE id = $1
	`, trabajoID)

	log.Printf("[embedding-worker] fragmento %s vectorizado con %s", fragmentoID, modelo)
}

func marcarFallido(ctx context.Context, trabajoID, fragmentoID string, intentos int, detalle string) {
	_, _ = db.Pool.Exec(ctx, `
		UPDATE trabajo_embedding
		SET estado = 'fallido', intentos = intentos + 1, error_detalle = $2
		WHERE id = $1
	`, trabajoID, detalle)

	_, _ = db.Pool.Exec(ctx, `
		UPDATE fragmento_conocimiento
		SET estado_embedding = 'error'
		WHERE id = $1
	`, fragmentoID)

	log.Printf("[embedding-worker] fragmento %s falló tras %d intentos: %s", fragmentoID, intentos+1, detalle)
}
