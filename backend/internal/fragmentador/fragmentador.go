package fragmentador

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"mi-tramite-bolivia-backend/internal/db"
)

// Fragmentar genera los fragmento_conocimiento para una versión recién publicada
// y crea los trabajo_embedding correspondientes. Debe llamarse dentro de la
// transacción de publicación o inmediatamente después.
func Fragmentar(ctx context.Context, versionID int64, modelo string) error {
	// Leer los datos de la versión
	var titulo, resumen, descripcion string
	var advertencias, plazoTexto *string
	err := db.Pool.QueryRow(ctx, `
		SELECT titulo, resumen, descripcion, advertencias, plazo_texto
		FROM tramite_version WHERE id = $1
	`, versionID).Scan(&titulo, &resumen, &descripcion, &advertencias, &plazoTexto)
	if err != nil {
		return fmt.Errorf("fragmentador: leer versión %d: %w", versionID, err)
	}

	// Fragmento de resumen
	partes := []string{
		"TRÁMITE: " + titulo,
		"RESUMEN: " + resumen,
		"DESCRIPCIÓN: " + descripcion,
	}
	if advertencias != nil && *advertencias != "" {
		partes = append(partes, "ADVERTENCIAS: "+*advertencias)
	}
	if plazoTexto != nil && *plazoTexto != "" {
		partes = append(partes, "PLAZO: "+*plazoTexto)
	}
	resumenContenido := strings.Join(partes, "\n\n")

	if err := insertarFragmento(ctx, versionID, "resumen", nil, resumenContenido, nil, modelo); err != nil {
		return err
	}

	// Fragmentos de requisitos
	rows, err := db.Pool.Query(ctx, `
		SELECT id, titulo, descripcion, aplica_si::text
		FROM requisito_tramite WHERE tramite_version_id = $1 ORDER BY orden
	`, versionID)
	if err != nil {
		return fmt.Errorf("fragmentador: leer requisitos: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var rid int64
		var rtitulo string
		var rdesc *string
		var raplicaSi string
		if err := rows.Scan(&rid, &rtitulo, &rdesc, &raplicaSi); err != nil {
			continue
		}
		lineas := []string{"REQUISITO: " + rtitulo}
		if rdesc != nil && *rdesc != "" {
			lineas = append(lineas, *rdesc)
		}
		if raplicaSi != "{}" && raplicaSi != "" {
			lineas = append(lineas, "CONDICIÓN: "+raplicaSi)
		}
		contenido := strings.Join(lineas, "\n")
		ridRef := rid
		if err := insertarFragmento(ctx, versionID, "requisito", &ridRef, contenido, nil, modelo); err != nil {
			return err
		}
	}

	// Fragmentos de pasos
	rowsPasos, err := db.Pool.Query(ctx, `
		SELECT id, numero, titulo, descripcion
		FROM paso_tramite WHERE tramite_version_id = $1 ORDER BY numero
	`, versionID)
	if err != nil {
		return fmt.Errorf("fragmentador: leer pasos: %w", err)
	}
	defer rowsPasos.Close()
	for rowsPasos.Next() {
		var pid int64
		var numero int
		var ptitulo, pdesc string
		if err := rowsPasos.Scan(&pid, &numero, &ptitulo, &pdesc); err != nil {
			continue
		}
		contenido := fmt.Sprintf("PASO %d: %s\n%s", numero, ptitulo, pdesc)
		pidRef := pid
		if err := insertarFragmento(ctx, versionID, "paso", &pidRef, contenido, nil, modelo); err != nil {
			return err
		}
	}

	// Fragmentos de costos
	rowsCostos, err := db.Pool.Query(ctx, `
		SELECT id, concepto, moneda, monto, monto_desde, monto_hasta, es_gratuito, medio_pago
		FROM costo_tramite WHERE tramite_version_id = $1
	`, versionID)
	if err != nil {
		return fmt.Errorf("fragmentador: leer costos: %w", err)
	}
	defer rowsCostos.Close()
	for rowsCostos.Next() {
		var cid int64
		var concepto, moneda string
		var monto, montoDesde, montoHasta *float64
		var esGratuito bool
		var medioPago *string
		if err := rowsCostos.Scan(&cid, &concepto, &moneda, &monto, &montoDesde, &montoHasta, &esGratuito, &medioPago); err != nil {
			continue
		}
		var costoStr string
		if esGratuito {
			costoStr = "GRATUITO"
		} else if monto != nil {
			costoStr = fmt.Sprintf("%.2f %s", *monto, moneda)
		} else if montoDesde != nil && montoHasta != nil {
			costoStr = fmt.Sprintf("%.2f - %.2f %s", *montoDesde, *montoHasta, moneda)
		} else {
			costoStr = "sin dato verificado"
		}
		contenido := fmt.Sprintf("COSTO: %s\nMonto: %s", concepto, costoStr)
		if medioPago != nil && *medioPago != "" {
			contenido += "\nMedio de pago: " + *medioPago
		}
		cidRef := cid
		if err := insertarFragmento(ctx, versionID, "costo", &cidRef, contenido, nil, modelo); err != nil {
			return err
		}
	}

	// Fragmentos de resultados
	rowsRes, err := db.Pool.Query(ctx, `
		SELECT id, nombre, vigencia, entrega
		FROM resultado_tramite WHERE tramite_version_id = $1 ORDER BY orden
	`, versionID)
	if err != nil {
		return fmt.Errorf("fragmentador: leer resultados: %w", err)
	}
	defer rowsRes.Close()
	for rowsRes.Next() {
		var rid int64
		var nombre string
		var vigencia, entrega *string
		if err := rowsRes.Scan(&rid, &nombre, &vigencia, &entrega); err != nil {
			continue
		}
		lineas := []string{"RESULTADO: " + nombre}
		if vigencia != nil && *vigencia != "" {
			lineas = append(lineas, "Vigencia: "+*vigencia)
		}
		if entrega != nil && *entrega != "" {
			lineas = append(lineas, "Entrega: "+*entrega)
		}
		contenido := strings.Join(lineas, "\n")
		ridRef := rid
		if err := insertarFragmento(ctx, versionID, "resultado", &ridRef, contenido, nil, modelo); err != nil {
			return err
		}
	}

	return nil
}

func insertarFragmento(
	ctx context.Context,
	versionID int64,
	tipo string,
	referenciaID *int64,
	contenido string,
	metadatos interface{},
	modelo string,
) error {
	hash := sha256hex(contenido)

	// Insertar o ignorar si ya existe (idempotente)
	var fragID string
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO fragmento_conocimiento
		    (tramite_version_id, tipo, referencia_id, contenido, metadatos,
		     hash_contenido, estado_embedding)
		VALUES ($1, $2, $3, $4, COALESCE($5::jsonb, '{}'::jsonb), $6, 'pendiente')
		ON CONFLICT (tramite_version_id, tipo, hash_contenido) DO NOTHING
		RETURNING id
	`, versionID, tipo, referenciaID, contenido, nil, hash).Scan(&fragID)

	if err != nil {
		// ON CONFLICT DO NOTHING devuelve 0 filas, no es error
		return nil
	}

	// Crear trabajo de embedding
	_, err = db.Pool.Exec(ctx, `
		INSERT INTO trabajo_embedding (fragmento_id, modelo)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, fragID, modelo)
	return err
}

func sha256hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
