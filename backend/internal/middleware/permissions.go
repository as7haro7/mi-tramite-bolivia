package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequirePermiso devuelve un middleware que verifica que el usuario tenga
// el permiso indicado. Los permisos usan el patrón "recurso:accion" o "recurso:*".
func RequirePermiso(permiso string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawPermisos, exists := c.Get("permisos")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    "PERMISSION_DENIED",
				"message": "Permisos no encontrados en el contexto",
			})
			return
		}

		permisos, ok := rawPermisos.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    "PERMISSION_DENIED",
				"message": "Formato de permisos inválido",
			})
			return
		}

		if !tienePermiso(permisos, permiso) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    "PERMISSION_DENIED",
				"message": "No tiene permiso para realizar esta acción",
			})
			return
		}

		c.Next()
	}
}

// tienePermiso verifica si la lista de permisos incluye el permiso solicitado.
// Soporta comodines: "usuarios:*" cubre cualquier "usuarios:xxx".
// También "*" global cubre todo.
func tienePermiso(permisos []string, requerido string) bool {
	partes := strings.SplitN(requerido, ":", 2)
	recursoReq := partes[0]
	accionReq := ""
	if len(partes) == 2 {
		accionReq = partes[1]
	}

	for _, p := range permisos {
		if p == "*" {
			return true
		}
		pp := strings.SplitN(p, ":", 2)
		if len(pp) != 2 {
			continue
		}
		recurso := pp[0]
		accion := pp[1]

		if recurso == recursoReq || recurso == "*" {
			if accion == "*" || accion == accionReq {
				return true
			}
		}
	}
	return false
}
