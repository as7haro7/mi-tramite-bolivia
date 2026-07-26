package admin

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"strings"
	"time"

	"mi-tramite-bolivia-backend/internal/db"
	"mi-tramite-bolivia-backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Correo   string `json:"correo" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Login godoc
// @Summary      Iniciar sesión administrativa
// @Description  Autentica a un usuario del panel y emite access + refresh token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Credenciales"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      401  {object} map[string]interface{}
// @Router       /api/v1/admin/auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	// Buscar usuario activo y sus roles
	var user db.UsuarioAdmin
	var passwordHash *string
	var roles []string

	err := db.Pool.QueryRow(context.Background(), `
		SELECT u.id, u.correo, u.nombre, u.password_hash, u.estado,
		       COALESCE(array_agg(r.codigo) FILTER (WHERE r.codigo IS NOT NULL), '{}') AS roles
		FROM usuario_admin u
		LEFT JOIN usuario_rol ur ON ur.usuario_id = u.id
		LEFT JOIN rol r ON r.id = ur.rol_id
		WHERE u.correo = $1 AND u.proveedor_identidad = 'local'
		GROUP BY u.id, u.correo, u.nombre, u.password_hash, u.estado
	`, strings.ToLower(req.Correo)).Scan(
		&user.ID, &user.Correo, &user.Nombre,
		&passwordHash, &user.Estado, &roles,
	)
	if err != nil {
		time.Sleep(200 * time.Millisecond) // timing blinding
		c.JSON(http.StatusUnauthorized, gin.H{"code": "AUTH_INVALID", "message": "Credenciales incorrectas"})
		return
	}

	if user.Estado != "activo" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "AUTH_INACTIVE", "message": "Cuenta inactiva o bloqueada"})
		return
	}

	if passwordHash == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "AUTH_NO_PASSWORD", "message": "Credenciales no configuradas. Use el flujo de activación."})
		return
	}

	// Verificar contraseña con bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(*passwordHash), []byte(req.Password)); err != nil {
		// Registrar intento fallido
		registrarIntentoFallido(context.Background(), user.ID, c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"code": "AUTH_INVALID", "message": "Credenciales incorrectas"})
		return
	}

	// Obtener permisos del usuario
	permisos, err := obtenerPermisos(context.Background(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "Error cargando permisos"})
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	accessTTL := 15 * time.Minute
	refreshTTL := 7 * 24 * time.Hour

	accessToken, err := middleware.GenerateAccessToken(jwtSecret, user.ID, user.Correo, permisos, accessTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "TOKEN_ERROR", "message": "Error generando token"})
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(jwtSecret, user.ID, user.Correo, refreshTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "TOKEN_ERROR", "message": "Error generando token"})
		return
	}

	// Guardar sesión con hash del refresh token
	refreshHash := hashSHA256(refreshToken)
	ipHash := hashSHA256(c.ClientIP())
	expira := time.Now().Add(refreshTTL)

	_, _ = db.Pool.Exec(context.Background(), `
		INSERT INTO sesion_admin
		    (usuario_id, refresh_token_hash, ip_hash, user_agent, expira_en)
		VALUES ($1, $2, $3, $4, $5)
	`, user.ID, refreshHash, ipHash, c.GetHeader("User-Agent"), expira)

	// Actualizar último acceso
	_, _ = db.Pool.Exec(context.Background(), `
		UPDATE usuario_admin SET ultimo_acceso_en = NOW() WHERE id = $1
	`, user.ID)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expira_en":     expira,
		"usuario": gin.H{
			"id":     user.ID,
			"correo": user.Correo,
			"nombre": user.Nombre,
			"roles":  roles,
		},
	})
}

// Refresh godoc
// @Summary      Renovar access token
// @Description  Usa el refresh token para emitir un nuevo par de tokens (rotación).
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body RefreshRequest true "Refresh token"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]interface{}
// @Router       /api/v1/admin/auth/refresh [post]
func Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	claims, err := middleware.ParseRefreshToken(jwtSecret, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "AUTH_INVALID", "message": "Refresh token inválido"})
		return
	}

	// Verificar que el token hash esté en la BD y no revocado
	refreshHash := hashSHA256(req.RefreshToken)
	var sesionID string
	err = db.Pool.QueryRow(context.Background(), `
		SELECT id FROM sesion_admin
		WHERE refresh_token_hash = $1
		  AND revocada_en IS NULL
		  AND expira_en > NOW()
	`, refreshHash).Scan(&sesionID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "AUTH_REVOKED", "message": "Sesión revocada o expirada"})
		return
	}

	// Verificar usuario activo
	var estado string
	err = db.Pool.QueryRow(context.Background(), `
		SELECT estado FROM usuario_admin WHERE id = $1
	`, claims.UsuarioID).Scan(&estado)
	if err != nil || estado != "activo" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "AUTH_INACTIVE", "message": "Usuario inactivo"})
		return
	}

	permisos, _ := obtenerPermisos(context.Background(), claims.UsuarioID)

	accessTTL := 15 * time.Minute
	refreshTTL := 7 * 24 * time.Hour

	newAccess, _ := middleware.GenerateAccessToken(jwtSecret, claims.UsuarioID, claims.Correo, permisos, accessTTL)
	newRefresh, _ := middleware.GenerateRefreshToken(jwtSecret, claims.UsuarioID, claims.Correo, refreshTTL)

	newHash := hashSHA256(newRefresh)
	expira := time.Now().Add(refreshTTL)

	// Revocar sesión anterior y crear nueva (rotación)
	_, _ = db.Pool.Exec(context.Background(), `
		UPDATE sesion_admin SET revocada_en = NOW() WHERE id = $1
	`, sesionID)

	_, _ = db.Pool.Exec(context.Background(), `
		INSERT INTO sesion_admin
		    (usuario_id, refresh_token_hash, ip_hash, user_agent, expira_en)
		VALUES ($1, $2, $3, $4, $5)
	`, claims.UsuarioID, newHash, hashSHA256(c.ClientIP()), c.GetHeader("User-Agent"), expira)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccess,
		"refresh_token": newRefresh,
		"expira_en":     expira,
	})
}

// Logout godoc
// @Summary      Cerrar sesión
// @Description  Revoca el refresh token del usuario autenticado.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body RefreshRequest false "Refresh token a revocar"
// @Success      200  {object} map[string]interface{}
// @Router       /api/v1/admin/auth/logout [post]
func Logout(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err == nil && req.RefreshToken != "" {
		hash := hashSHA256(req.RefreshToken)
		_, _ = db.Pool.Exec(context.Background(), `
			UPDATE sesion_admin SET revocada_en = NOW()
			WHERE refresh_token_hash = $1
		`, hash)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada"})
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func obtenerPermisos(ctx context.Context, usuarioID int64) ([]string, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT DISTINCT p.value
		FROM usuario_rol ur
		JOIN rol r ON r.id = ur.rol_id
		CROSS JOIN LATERAL jsonb_array_elements_text(r.permisos) AS p(value)
		WHERE ur.usuario_id = $1
	`, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permisos []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err == nil {
			permisos = append(permisos, p)
		}
	}
	return permisos, nil
}

func registrarIntentoFallido(ctx context.Context, usuarioID int64, ip string) {
	// Contar intentos recientes; bloquear si supera umbral
	var intentos int
	_ = db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM evento_auditoria
		WHERE actor_id = $1
		  AND accion = 'login.fallido'
		  AND ocurrido_en > NOW() - INTERVAL '15 minutes'
	`, usuarioID).Scan(&intentos)

	if intentos >= 5 {
		_, _ = db.Pool.Exec(ctx, `
			UPDATE usuario_admin SET estado = 'bloqueado' WHERE id = $1 AND estado = 'activo'
		`, usuarioID)
	}

	ipHash := hashSHA256(ip)
	idStr := int64ToStr(usuarioID)
	_, _ = db.Pool.Exec(ctx, `
		INSERT INTO evento_auditoria
		    (actor_id, accion, entidad_tipo, entidad_id, ip_hash)
		VALUES ($1, 'login.fallido', 'usuario_admin', $2, $3)
	`, usuarioID, idStr, ipHash)
}

func hashSHA256(s string) string {
	h := hmac.New(sha256.New, []byte(os.Getenv("IP_HASH_SECRET")+"salt"))
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func int64ToStr(n int64) string {
	b := make([]byte, 0, 20)
	if n < 0 {
		b = append(b, '-')
		n = -n
	}
	if n == 0 {
		return "0"
	}
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}
