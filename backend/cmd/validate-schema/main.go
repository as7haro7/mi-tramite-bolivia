// Command validate-schema executes the development schema in an isolated
// PostgreSQL schema and always rolls the transaction back. It validates the
// complete DDL and seed without touching the application's public tables.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

const (
	resetMarker   = "-- -----------------------------------------------------------------------------\n-- RESET EXPLÍCITO"
	catalogMarker = "-- -----------------------------------------------------------------------------\n-- CATÁLOGOS TERRITORIALES"
)

func main() {
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL no está configurada")
	}

	seedPath := "../mi-tramite-bolivia-seed.sql"
	if len(os.Args) > 1 {
		seedPath = os.Args[1]
	}

	content, err := os.ReadFile(seedPath)
	if err != nil {
		log.Fatalf("no se pudo leer el SQL: %v", err)
	}

	validationSchema := fmt.Sprintf("mtb_validation_%d", time.Now().UnixNano())
	sql, err := isolateSQL(string(content), validationSchema)
	if err != nil {
		log.Fatal(err)
	}

	config, err := pgx.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("DATABASE_URL inválida: %v", err)
	}
	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		log.Fatalf("no se pudo conectar con PostgreSQL: %v", err)
	}
	defer conn.Close(context.Background())

	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatalf("no se pudo iniciar la transacción: %v", err)
	}
	defer tx.Rollback(context.Background())

	if _, err = tx.Exec(ctx, "CREATE SCHEMA "+pgx.Identifier{validationSchema}.Sanitize()); err != nil {
		log.Fatalf("no se pudo crear el esquema aislado: %v", err)
	}

	if _, err = tx.Exec(ctx, sql); err != nil {
		log.Fatalf("el esquema o la semilla no son válidos: %v", err)
	}

	if err = tx.Rollback(ctx); err != nil {
		log.Fatalf("la validación terminó, pero falló el rollback: %v", err)
	}

	fmt.Println("Esquema y datos semilla válidos; la transacción de prueba fue revertida.")
}

func isolateSQL(input, schema string) (string, error) {
	sql := strings.ReplaceAll(input, "\r\n", "\n")

	resetStart := strings.Index(sql, resetMarker)
	catalogStart := strings.Index(sql, catalogMarker)
	if resetStart == -1 || catalogStart == -1 || catalogStart <= resetStart {
		return "", fmt.Errorf("no se encontraron los marcadores del bloque RESET")
	}

	// El reset se omite para que la búsqueda de nombres nunca alcance tablas
	// existentes fuera del esquema de validación.
	sql = sql[:resetStart] + sql[catalogStart:]
	sql = strings.Replace(sql, "BEGIN;\n", "", 1)
	sql = strings.Replace(sql, "\nCOMMIT;\n", "\n", 1)

	searchPath := fmt.Sprintf(
		"SET search_path TO %s, public, extensions, pg_catalog;",
		pgx.Identifier{schema}.Sanitize(),
	)
	sql = strings.ReplaceAll(
		sql,
		"SET search_path TO public, extensions, pg_catalog;",
		searchPath,
	)
	sql = strings.ReplaceAll(
		sql,
		"public.normalizar_busqueda",
		pgx.Identifier{schema}.Sanitize()+".normalizar_busqueda",
	)

	if strings.Contains(sql, "\nCOMMIT;") || strings.Contains(sql, "\nBEGIN;") {
		return "", fmt.Errorf("el SQL aislado todavía contiene control de transacción")
	}

	return sql, nil
}
