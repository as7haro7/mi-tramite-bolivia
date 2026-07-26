package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"mi-tramite-bolivia-backend/internal/config"
	"mi-tramite-bolivia-backend/internal/db"
	"mi-tramite-bolivia-backend/internal/handlers"
	admhandlers "mi-tramite-bolivia-backend/internal/handlers/admin"
	"mi-tramite-bolivia-backend/internal/middleware"
	"mi-tramite-bolivia-backend/internal/worker"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "mi-tramite-bolivia-backend/docs"
)

// @title           Mi Trámite Bolivia API v2
// @version         2.0
// @description     API para orientación sobre trámites públicos en Bolivia. Información oficial versionada y trazable.
// @contact.name    Equipo Mi Trámite Bolivia

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT access token obtenido mediante POST /api/v1/admin/auth/login

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("[config] no se encontró .env, usando variables de entorno del sistema")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[config] configuración inválida: %v", err)
	}

	if err := db.ConnectDB(); err != nil {
		log.Fatalf("[db] error conectando: %v", err)
	}
	defer db.Pool.Close()

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS — en producción solo orígenes de la config; en desarrollo permite todos
	var corsConfig cors.Config
	if cfg.Environment == "production" {
		corsConfig = cors.Config{
			AllowOrigins:     cfg.CORSOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
			ExposeHeaders:    []string{"X-Request-ID"},
			AllowCredentials: false,
			MaxAge:           12 * time.Hour,
		}
	} else {
		corsConfig = cors.DefaultConfig()
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"}
	}
	r.Use(cors.New(corsConfig))

	r.Use(middleware.RateLimitMiddleware())

	// ─── Health checks ──────────────────────────────────────────────────────────
	r.GET("/health/live", handlers.HealthLive)
	r.GET("/health/ready", handlers.HealthReady)

	// ─── Swagger ────────────────────────────────────────────────────────────────
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ─── API v1 pública ─────────────────────────────────────────────────────────
	v1 := r.Group("/api/v1")
	{
		// Catálogo
		v1.GET("/tramites", handlers.GetTramites)
		v1.GET("/tramites/:slug", handlers.GetTramiteBySlug)
		v1.GET("/tramites/:slug/oficinas", handlers.GetOficinasDeModalidad)
		v1.GET("/instituciones", handlers.GetInstituciones)
		v1.GET("/categorias", handlers.GetCategorias)
		v1.GET("/oficinas", handlers.GetOficinas)

		// Chat RAG
		v1.POST("/chat/conversaciones", handlers.CrearConversacion)
		v1.POST("/chat/conversaciones/:id/mensajes", handlers.EnviarMensaje)
		v1.POST("/mensajes/:id/feedback", handlers.RegistrarFeedback)

		// Reportes ciudadanos (anónimo)
		v1.POST("/reportes", handlers.CrearReporte)
	}

	// ─── API v1 administrativa ──────────────────────────────────────────────────
	jwtMW := middleware.JWTMiddleware(cfg.JWTSecret)
	perm := middleware.RequirePermiso

	admin := v1.Group("/admin")
	{
		// Auth (sin JWT)
		admin.POST("/auth/login", admhandlers.Login)
		admin.POST("/auth/refresh", admhandlers.Refresh)
		admin.POST("/auth/logout", admhandlers.Logout)

		// Dashboard (requiere autenticación)
		admin.GET("/dashboard", jwtMW, admhandlers.Dashboard)

		// Trámites y versiones
		tramitesAdm := admin.Group("/tramites", jwtMW)
		{
			tramitesAdm.GET("", admhandlers.ListarTramitesAdmin)
			tramitesAdm.GET("/:id", admhandlers.ObtenerTramiteAdmin)
			tramitesAdm.POST("", perm("catalogo:*"), admhandlers.CrearTramiteAdmin)
			tramitesAdm.PUT("/:id", perm("catalogo:*"), admhandlers.ActualizarTramiteAdmin)
			tramitesAdm.DELETE("/:id", perm("catalogo:*"), admhandlers.EliminarTramiteAdmin)
			tramitesAdm.POST("/:id/versiones", perm("version:crear"), admhandlers.CrearVersion)
		}

		versionesAdm := admin.Group("/versiones", jwtMW)
		{
			versionesAdm.PUT("/:id", perm("version:editar"), admhandlers.ActualizarVersion)
			versionesAdm.POST("/:id/enviar-revision", perm("version:editar"), admhandlers.EnviarARevision)
			versionesAdm.POST("/:id/aprobar", perm("version:aprobar"), admhandlers.AprobarVersion)
			versionesAdm.POST("/:id/rechazar", perm("version:aprobar"), admhandlers.RechazarVersion)
			versionesAdm.POST("/:id/publicar", perm("publicacion:*"), admhandlers.PublicarVersion)
		}

		// Categorias
		categoriasAdm := admin.Group("/categorias", jwtMW)
		{
			categoriasAdm.GET("", admhandlers.ListarCategoriasAdmin)
			categoriasAdm.POST("", perm("catalogo:*"), admhandlers.CrearCategoria)
			categoriasAdm.PUT("/:id", perm("catalogo:*"), admhandlers.ActualizarCategoria)
			categoriasAdm.DELETE("/:id", perm("catalogo:*"), admhandlers.EliminarCategoria)
		}

		// Etiquetas
		etiquetasAdm := admin.Group("/etiquetas", jwtMW)
		{
			etiquetasAdm.GET("", admhandlers.ListarEtiquetasAdmin)
			etiquetasAdm.POST("", perm("catalogo:*"), admhandlers.CrearEtiqueta)
			etiquetasAdm.PUT("/:id", perm("catalogo:*"), admhandlers.ActualizarEtiqueta)
			etiquetasAdm.DELETE("/:id", perm("catalogo:*"), admhandlers.EliminarEtiqueta)
		}

		// Instituciones y oficinas
		instAdm := admin.Group("/instituciones", jwtMW)
		{
			instAdm.GET("", admhandlers.ListarInstitucionesAdmin)
			instAdm.POST("", perm("catalogo:*"), admhandlers.CrearInstitucionAdmin)
			instAdm.PUT("/:id", perm("catalogo:*"), admhandlers.ActualizarInstitucionAdmin)
			instAdm.DELETE("/:id", perm("catalogo:*"), admhandlers.EliminarInstitucionAdmin)
		}

		oficinasAdm := admin.Group("/oficinas", jwtMW)
		{
			oficinasAdm.GET("", admhandlers.ListarOficinasAdmin)
			oficinasAdm.POST("", perm("catalogo:*"), admhandlers.CrearOficinaAdmin)
			oficinasAdm.PUT("/:id", perm("catalogo:*"), admhandlers.ActualizarOficinaAdmin)
			oficinasAdm.DELETE("/:id", perm("catalogo:*"), admhandlers.EliminarOficinaAdmin)
		}

		// Ingesta
		admin.GET("/fuentes", jwtMW, admhandlers.ListarFuentes)
		admin.POST("/fuentes", jwtMW, perm("ingesta:*"), admhandlers.CrearFuente)
		admin.PUT("/fuentes/:id", jwtMW, perm("ingesta:*"), admhandlers.ActualizarFuente)
		admin.DELETE("/fuentes/:id", jwtMW, perm("ingesta:*"), admhandlers.EliminarFuente)
		admin.GET("/ingestas", jwtMW, admhandlers.ListarIngestas)
		admin.GET("/candidatos", jwtMW, admhandlers.ListarCandidatos)
		admin.PUT("/candidatos/:id", jwtMW, perm("ingesta:*"), admhandlers.ActualizarCandidato)

		// Usuarios
		usuariosAdm := admin.Group("/usuarios", jwtMW, perm("usuarios:*"))
		{
			usuariosAdm.GET("", admhandlers.ListarUsuarios)
			usuariosAdm.POST("", admhandlers.CrearUsuario)
			usuariosAdm.PUT("/:id", admhandlers.ActualizarUsuario)
		}

		// Auditoría (solo lectura)
		admin.GET("/auditoria", jwtMW, perm("auditoria:leer"), admhandlers.ListarAuditoria)

		// Reportes ciudadanos
		reportesAdm := admin.Group("/reportes", jwtMW)
		{
			reportesAdm.GET("", admhandlers.ListarReportesAdmin)
			reportesAdm.PUT("/:id", perm("reporte:gestionar"), admhandlers.ActualizarReporte)
		}
	}

	// ─── Worker de embeddings ────────────────────────────────────────────────────
	workerCtx, cancelWorker := context.WithCancel(context.Background())
	go worker.StartEmbeddingWorker(workerCtx)

	// ─── Servidor HTTP con graceful shutdown ─────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("[server] escuchando en http://localhost:%s", cfg.Port)
		log.Printf("[server] ambiente: %s", cfg.Environment)
		log.Printf("[server] CORS origins: %s", strings.Join(cfg.CORSOrigins, ", "))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[server] error: %v", err)
		}
	}()

	// Esperar señal de terminación
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[server] apagando...")
	cancelWorker()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[server] error en shutdown: %v", err)
	}
	log.Println("[server] apagado limpiamente")
}
