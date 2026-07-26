package admin

import (
	"context"
	"fmt"
	"net/http"

	"mi-tramite-bolivia-backend/internal/auditoria"
	"mi-tramite-bolivia-backend/internal/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ListarUsuarios godoc
// @Summary Listar usuarios administrativos
// @Tags    Admin - Usuarios
// @Produce json
// @Security BearerAuth
// @Router  /api/v1/admin/usuarios [get]
func ListarUsuarios(c *gin.Context) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT u.id, u.correo, u.nombre, u.estado,
		       u.ultimo_acceso_en, u.creado_en,
		       array_agg(r.codigo) FILTER (WHERE r.codigo IS NOT NULL) AS roles
		FROM usuario_admin u
		LEFT JOIN usuario_rol ur ON ur.usuario_id = u.id
		LEFT JOIN rol r ON r.id = ur.rol_id
		GROUP BY u.id, u.correo, u.nombre, u.estado, u.ultimo_acceso_en, u.creado_en
		ORDER BY u.creado_en DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}
	defer rows.Close()

	type UItem struct {
		db.UsuarioAdmin
		Roles []string `json:"roles"`
	}
	datos := make([]UItem, 0)
	for rows.Next() {
		var u UItem
		if err := rows.Scan(
			&u.ID, &u.Correo, &u.Nombre, &u.Estado,
			&u.UltimoAccesoEn, &u.CreadoEn,
			&u.Roles,
		); err != nil {
			continue
		}
		if u.Roles == nil {
			u.Roles = []string{}
		}
		datos = append(datos, u)
	}

	c.JSON(http.StatusOK, gin.H{"datos": datos})
}

// CrearUsuario godoc
// @Summary Crear usuario administrativo
// @Tags    Admin - Usuarios
// @Accept  json
// @Produce json
// @Security BearerAuth
// @Router  /api/v1/admin/usuarios [post]
func CrearUsuario(c *gin.Context) {
	var req struct {
		Correo   string `json:"correo" binding:"required,email"`
		Nombre   string `json:"nombre" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
		RolCodigo string `json:"rol_codigo" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "HASH_ERROR", "message": "Error procesando contraseña"})
		return
	}

	var userID int64
	err = db.Pool.QueryRow(context.Background(), `
		INSERT INTO usuario_admin (correo, nombre, password_hash, estado)
		VALUES ($1, $2, $3, 'activo') RETURNING id
	`, req.Correo, req.Nombre, string(hash)).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error creando usuario: " + err.Error()})
		return
	}

	// Asignar rol
	var rolID int
	if err := db.Pool.QueryRow(context.Background(), `SELECT id FROM rol WHERE codigo = $1`, req.RolCodigo).Scan(&rolID); err != nil {
		// Crear sin rol si no existe
		c.JSON(http.StatusCreated, gin.H{"id": userID, "warning": "rol no encontrado"})
		return
	}

	_, _ = db.Pool.Exec(context.Background(), `
		INSERT INTO usuario_rol (usuario_id, rol_id) VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, userID, rolID)

	actorID := actorIDFromCtx(c)
	idStr := fmt.Sprintf("%d", userID)
	auditoria.Registrar(c.Request.Context(), actorID, "usuario.crear", "usuario_admin", &idStr, nil, map[string]string{"correo": req.Correo, "rol": req.RolCodigo}, nil, nil)
	c.JSON(http.StatusCreated, gin.H{"id": userID})
}

// ActualizarUsuario godoc
// @Summary Actualizar estado/rol de usuario
// @Tags    Admin - Usuarios
// @Accept  json
// @Produce json
// @Security BearerAuth
// @Param   id path int true "ID del usuario"
// @Router  /api/v1/admin/usuarios/{id} [put]
func ActualizarUsuario(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Estado    *string `json:"estado"`
		Nombre    *string `json:"nombre"`
		RolCodigo *string `json:"rol_codigo"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	var antes db.UsuarioAdmin
	_ = db.Pool.QueryRow(context.Background(), `SELECT id, estado, nombre FROM usuario_admin WHERE id = $1`, id).
		Scan(&antes.ID, &antes.Estado, &antes.Nombre)

	_, err := db.Pool.Exec(context.Background(), `
		UPDATE usuario_admin
		SET estado = COALESCE($2, estado),
		    nombre = COALESCE($3, nombre)
		WHERE id = $1
	`, id, req.Estado, req.Nombre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": err.Error()})
		return
	}

	if req.RolCodigo != nil {
		var rolID int
		if err := db.Pool.QueryRow(context.Background(), `SELECT id FROM rol WHERE codigo = $1`, *req.RolCodigo).Scan(&rolID); err == nil {
			_, _ = db.Pool.Exec(context.Background(), `DELETE FROM usuario_rol WHERE usuario_id = $1`, id)
			_, _ = db.Pool.Exec(context.Background(), `INSERT INTO usuario_rol (usuario_id, rol_id) VALUES ($1, $2)`, id, rolID)
		}
	}

	// Revocar sesiones si se bloquea o cambia rol
	if (req.Estado != nil && (*req.Estado == "bloqueado" || *req.Estado == "revocado")) || req.RolCodigo != nil {
		_, _ = db.Pool.Exec(context.Background(), `
			UPDATE sesion_admin SET revocada_en = NOW()
			WHERE usuario_id = $1 AND revocada_en IS NULL
		`, id)
	}

	actorID := actorIDFromCtx(c)
	auditoria.Registrar(c.Request.Context(), actorID, "usuario.actualizar", "usuario_admin", &id, antes, req, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado"})
}
