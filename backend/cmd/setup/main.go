// cmd/setup/main.go
// Script de configuración inicial: crea el primer usuario administrador.
// Uso: go run ./cmd/setup/main.go
//
// Variables de entorno requeridas: DATABASE_URL, ADMIN_CORREO, ADMIN_PASSWORD, ADMIN_NOMBRE

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("[setup] no se encontró .env")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("[setup] DATABASE_URL no configurado")
	}

	correo := os.Getenv("ADMIN_CORREO")
	password := os.Getenv("ADMIN_PASSWORD")
	nombre := os.Getenv("ADMIN_NOMBRE")

	if correo == "" {
		correo = "admin@mitramite.bo"
	}
	if nombre == "" {
		nombre = "Administrador Principal"
	}
	if password == "" {
		log.Fatal("[setup] ADMIN_PASSWORD no configurado")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("[setup] error conectando a BD: %v", err)
	}
	defer pool.Close()

	// Verificar que la tabla existe
	var exists bool
	_ = pool.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = 'usuario_admin'
		)
	`).Scan(&exists)

	if !exists {
		log.Fatal("[setup] la tabla usuario_admin no existe. ¿El seed fue ejecutado?")
	}

	// Hash de la contraseña
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("[setup] error generando hash: %v", err)
	}

	// Crear usuario
	var userID int64
	err = pool.QueryRow(context.Background(), `
		INSERT INTO usuario_admin (correo, nombre, password_hash, proveedor_identidad, estado)
		VALUES ($1, $2, $3, 'local', 'activo')
		ON CONFLICT (correo) DO UPDATE
		    SET password_hash = EXCLUDED.password_hash,
		        estado = 'activo'
		RETURNING id
	`, correo, nombre, string(hash)).Scan(&userID)
	if err != nil {
		log.Fatalf("[setup] error creando usuario: %v", err)
	}

	// Asignar rol superadmin si existe
	var rolID int
	err = pool.QueryRow(context.Background(), `
		SELECT id FROM rol WHERE codigo = 'superadmin' LIMIT 1
	`).Scan(&rolID)
	if err == nil {
		_, _ = pool.Exec(context.Background(), `
			INSERT INTO usuario_rol (usuario_id, rol_id) VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, userID, rolID)
		fmt.Printf("[setup] rol superadmin asignado\n")
	} else {
		// Crear rol superadmin si no existe
		var newRolID int
		err = pool.QueryRow(context.Background(), `
			INSERT INTO rol (codigo, nombre, permisos)
			VALUES ('superadmin', 'Superadministrador', '["*"]'::jsonb)
			ON CONFLICT (codigo) DO UPDATE SET permisos = EXCLUDED.permisos
			RETURNING id
		`).Scan(&newRolID)
		if err == nil {
			_, _ = pool.Exec(context.Background(), `
				INSERT INTO usuario_rol (usuario_id, rol_id) VALUES ($1, $2)
				ON CONFLICT DO NOTHING
			`, userID, newRolID)
			fmt.Printf("[setup] rol superadmin creado y asignado\n")
		}
	}

	fmt.Printf("[setup] ✓ Usuario administrador creado:\n")
	fmt.Printf("  ID:     %d\n", userID)
	fmt.Printf("  Correo: %s\n", correo)
	fmt.Printf("  Nombre: %s\n", nombre)
	fmt.Printf("  Rol:    superadmin\n")
	fmt.Printf("\nAhora puede iniciar sesión en:\n  POST /api/v1/admin/auth/login\n")
}
