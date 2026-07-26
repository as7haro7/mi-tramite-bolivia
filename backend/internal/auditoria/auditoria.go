package auditoria

import (
	"context"
	"encoding/json"
	"log"

	"mi-tramite-bolivia-backend/internal/db"
)

// Registrar inserta un evento de auditoría de forma no bloqueante.
// No registra contraseñas, tokens ni datos sensibles.
func Registrar(
	ctx context.Context,
	actorID *int64,
	accion string,
	entidadTipo string,
	entidadID *string,
	antes interface{},
	despues interface{},
	ipHash *string,
	userAgent *string,
) {
	var antesJSON, despuesJSON []byte
	if antes != nil {
		antesJSON, _ = json.Marshal(antes)
	}
	if despues != nil {
		despuesJSON, _ = json.Marshal(despues)
	}

	_, err := db.Pool.Exec(ctx, `
		INSERT INTO evento_auditoria
		    (actor_id, accion, entidad_tipo, entidad_id, antes, despues, ip_hash, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, actorID, accion, entidadTipo, entidadID,
		nullableJSON(antesJSON), nullableJSON(despuesJSON),
		ipHash, userAgent,
	)
	if err != nil {
		log.Printf("[auditoria] error registrando evento %s/%s: %v", entidadTipo, accion, err)
	}
}

func nullableJSON(b []byte) interface{} {
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	return b
}
