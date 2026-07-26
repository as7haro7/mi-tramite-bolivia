package middleware

import (
	"net/http"
	"strings"
	"time"

	"mi-tramite-bolivia-backend/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware valida el access token JWT y carga los claims en el contexto.
func JWTMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "AUTH_REQUIRED",
				"message": "Se requiere autorización",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "AUTH_INVALID_FORMAT",
				"message": "Formato inválido, se espera 'Bearer <token>'",
			})
			return
		}

		tokenStr := parts[1]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "AUTH_TOKEN_INVALID",
				"message": "Token inválido o expirado",
			})
			return
		}

		if claims.Type != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "AUTH_WRONG_TOKEN_TYPE",
				"message": "Se requiere un access token",
			})
			return
		}

		// Verificar que el usuario todavía esté activo
		var estado string
		err = db.Pool.QueryRow(c.Request.Context(),
			`SELECT estado FROM usuario_admin WHERE id = $1`, claims.UsuarioID,
		).Scan(&estado)
		if err != nil || estado != "activo" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "AUTH_USER_INACTIVE",
				"message": "Usuario inactivo o no encontrado",
			})
			return
		}

		c.Set("usuario_id", claims.UsuarioID)
		c.Set("correo", claims.Correo)
		c.Set("permisos", claims.Permisos)
		c.Next()
	}
}

// Claims es el payload del JWT de administración.
type Claims struct {
	UsuarioID int64    `json:"usuario_id"`
	Correo    string   `json:"correo"`
	Permisos  []string `json:"permisos"`
	Type      string   `json:"type"` // "access" | "refresh"
	jwt.RegisteredClaims
}

// GenerateAccessToken emite un access token de vida corta.
func GenerateAccessToken(secret string, usuarioID int64, correo string, permisos []string, ttl time.Duration) (string, error) {
	claims := Claims{
		UsuarioID: usuarioID,
		Correo:    correo,
		Permisos:  permisos,
		Type:      "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken emite un refresh token opaco largo.
func GenerateRefreshToken(secret string, usuarioID int64, correo string, ttl time.Duration) (string, error) {
	claims := Claims{
		UsuarioID: usuarioID,
		Correo:    correo,
		Type:      "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseRefreshToken parsea y valida un refresh token devolviendo los claims.
func ParseRefreshToken(secret, tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid || claims.Type != "refresh" {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
