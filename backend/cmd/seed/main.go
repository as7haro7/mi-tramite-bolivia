package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"mi-tramite-bolivia-backend/internal/db"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: No se encontró archivo .env local")
	}

	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	defer db.Pool.Close()

	ctx := context.Background()

	instituciones := []struct{ nombre, sigla, url string }{
		{"Servicio General de Identificación Personal", "SEGIP", "https://segip.gob.bo"},
		{"Servicio de Registro Cívico", "SERECI", "https://sereci.gob.bo"},
		{"Organismo Operativo de Tránsito", "Tránsito", "https://transito.gob.bo"},
	}

	for _, inst := range instituciones {
		_, err := db.Pool.Exec(ctx, `
			INSERT INTO institucion (nombre, sigla, url_oficial) 
			VALUES ($1, $2, $3) 
			ON CONFLICT (sigla) DO NOTHING
		`, inst.nombre, inst.sigla, inst.url)
		if err != nil {
			log.Printf("Error insertando %s: %v", inst.sigla, err)
		} else {
			fmt.Printf("Institución %s configurada.\n", inst.sigla)
		}
	}

	username := "as7haro7"
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		log.Println("ADVERTENCIA: No se definió ADMIN_PASSWORD en el .env, usando 'admin123' temporalmente.")
		password = "admin123"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error generando hash: %v", err)
	}

	_, err = db.Pool.Exec(ctx, `
		INSERT INTO administrador (nombre_usuario, password_hash, rol) 
		VALUES ($1, $2, 'admin') 
		ON CONFLICT (nombre_usuario) DO UPDATE SET password_hash = $2
	`, username, string(hash))

	if err != nil {
		log.Fatalf("Error creando administrador: %v", err)
	}

	fmt.Println("---------------------------------------------------")
	fmt.Println("Datos iniciales (Seed) creados exitosamente.")
	fmt.Println("Credenciales de acceso para el Panel Web:")
	fmt.Printf("Usuario: %s\n", username)
	fmt.Printf("Contraseña: %s\n", password)
	fmt.Println("---------------------------------------------------")
}
