package main

import (
	"context"
	"log"

	"mi-tramite-bolivia-backend/internal/config"
	"mi-tramite-bolivia-backend/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env found, using system env")
	}

	_, err := config.Load()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	if err := db.ConnectDB(); err != nil {
		log.Fatalf("DB error: %v", err)
	}
	defer db.Pool.Close()

	res, err := db.Pool.Exec(context.Background(), `UPDATE rol SET permisos = '["*"]' WHERE codigo = 'superadmin'`)
	if err != nil {
		log.Fatalf("Query error: %v", err)
	}

	log.Printf("Roles actualizados: %v", res.RowsAffected())
}
