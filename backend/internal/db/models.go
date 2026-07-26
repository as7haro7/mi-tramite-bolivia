package db

import (
	"encoding/json"
	"time"
)

// ─── Territorios ───────────────────────────────────────────────────────────────

type Departamento struct {
	ID     int    `json:"id"`
	Codigo string `json:"codigo"`
	Nombre string `json:"nombre"`
	Activo bool   `json:"activo"`
}

type Municipio struct {
	ID             int64   `json:"id"`
	DepartamentoID int     `json:"departamento_id"`
	CodigoINE      *string `json:"codigo_ine,omitempty"`
	Nombre         string  `json:"nombre"`
	Activo         bool    `json:"activo"`
}

// ─── Multimedia ────────────────────────────────────────────────────────────────

type ArchivoMultimedia struct {
	ID               string    `json:"id"`
	OriginalID       *string   `json:"original_id,omitempty"`
	Variante         string    `json:"variante"`
	Proveedor        string    `json:"proveedor"`
	Contenedor       string    `json:"contenedor"`
	ClaveObjeto      string    `json:"clave_objeto"`
	URLPublica       *string   `json:"url_publica,omitempty"`
	MimeType         string    `json:"mime_type"`
	TamanoBytes      int64     `json:"tamano_bytes"`
	AnchoPx          *int      `json:"ancho_px,omitempty"`
	AltoPx           *int      `json:"alto_px,omitempty"`
	SHA256           string    `json:"sha256"`
	TextoAlternativo string    `json:"texto_alternativo"`
	Estado           string    `json:"estado"`
	CreadoPor        *int64    `json:"creado_por,omitempty"`
	CreadoEn         time.Time `json:"creado_en"`
	ActualizadoEn    time.Time `json:"actualizado_en"`
}

// ─── Instituciones ─────────────────────────────────────────────────────────────

type Institucion struct {
	ID                int64     `json:"id"`
	Codigo            string    `json:"codigo"`
	Nombre            string    `json:"nombre"`
	Sigla             string    `json:"sigla"`
	Tipo              string    `json:"tipo"`
	InstitucionPadreID *int64   `json:"institucion_padre_id,omitempty"`
	SitioWeb          *string   `json:"sitio_web,omitempty"`
	TelefonoGeneral   *string   `json:"telefono_general,omitempty"`
	CorreoGeneral     *string   `json:"correo_general,omitempty"`
	LogoArchivoID     *string   `json:"logo_archivo_id,omitempty"`
	Estado            string    `json:"estado"`
	ReemplazadaPorID  *int64    `json:"reemplazada_por_id,omitempty"`
	CreadoEn          time.Time `json:"creado_en"`
	ActualizadoEn     time.Time `json:"actualizado_en"`
}

type InstitucionAlias struct {
	ID            int64  `json:"id"`
	InstitucionID int64  `json:"institucion_id"`
	Alias         string `json:"alias"`
	Normalizado   string `json:"normalizado"`
}

// ─── Oficinas ──────────────────────────────────────────────────────────────────

type Oficina struct {
	ID            int64      `json:"id"`
	InstitucionID int64      `json:"institucion_id"`
	MunicipioID   *int64     `json:"municipio_id,omitempty"`
	Codigo        string     `json:"codigo"`
	Nombre        string     `json:"nombre"`
	Tipo          string     `json:"tipo"`
	Direccion     *string    `json:"direccion,omitempty"`
	Referencia    *string    `json:"referencia,omitempty"`
	Latitud       *float64   `json:"latitud,omitempty"`
	Longitud      *float64   `json:"longitud,omitempty"`
	ZonaHoraria   string     `json:"zona_horaria"`
	Accesible     *bool      `json:"accesible,omitempty"`
	RequiereCita  bool       `json:"requiere_cita"`
	URLCita       *string    `json:"url_cita,omitempty"`
	Estado        string     `json:"estado"`
	VerificadoEn  *time.Time `json:"verificado_en,omitempty"`
	CreadoEn      time.Time  `json:"creado_en"`
	ActualizadoEn time.Time  `json:"actualizado_en"`
}

type OficinaSummary struct {
	ID           int64    `json:"id"`
	Nombre       string   `json:"nombre"`
	Tipo         string   `json:"tipo"`
	Direccion    *string  `json:"direccion,omitempty"`
	Latitud      *float64 `json:"latitud,omitempty"`
	Longitud     *float64 `json:"longitud,omitempty"`
	RequiereCita bool     `json:"requiere_cita"`
	URLCita      *string  `json:"url_cita,omitempty"`
	Estado       string   `json:"estado"`
	Municipio    *string  `json:"municipio,omitempty"`
	Departamento *string  `json:"departamento,omitempty"`
}

type HorarioOficina struct {
	ID           int64   `json:"id"`
	OficinaID    int64   `json:"oficina_id"`
	DiaSemana    int     `json:"dia_semana"`
	HoraApertura string  `json:"hora_apertura"`
	HoraCierre   string  `json:"hora_cierre"`
	ConCita      bool    `json:"con_cita"`
	ValidoDesde  *string `json:"valido_desde,omitempty"`
	ValidoHasta  *string `json:"valido_hasta,omitempty"`
}

type ExcepcionAtencion struct {
	ID           int64   `json:"id"`
	OficinaID    int64   `json:"oficina_id"`
	Fecha        string  `json:"fecha"`
	Cerrada      bool    `json:"cerrada"`
	HoraApertura *string `json:"hora_apertura,omitempty"`
	HoraCierre   *string `json:"hora_cierre,omitempty"`
	Motivo       *string `json:"motivo,omitempty"`
}

// ─── Catálogo de trámites ──────────────────────────────────────────────────────

type Categoria struct {
	ID      int64   `json:"id"`
	PadreID *int64  `json:"padre_id,omitempty"`
	Codigo  string  `json:"codigo"`
	Nombre  string  `json:"nombre"`
	Icono   *string `json:"icono,omitempty"`
	Orden   int     `json:"orden"`
	Activa  bool    `json:"activa"`
}

type Tramite struct {
	ID                  int64     `json:"id"`
	Codigo              string    `json:"codigo"`
	Slug                string    `json:"slug"`
	InstitucionID       int64     `json:"institucion_id"`
	CategoriaID         *int64    `json:"categoria_id,omitempty"`
	CodigoOficial       *string   `json:"codigo_oficial,omitempty"`
	Alcance             string    `json:"alcance"`
	Estado              string    `json:"estado"`
	ReemplazaTramiteID  *int64    `json:"reemplaza_tramite_id,omitempty"`
	CreadoEn            time.Time `json:"creado_en"`
	ActualizadoEn       time.Time `json:"actualizado_en"`
}

type TramiteAlias struct {
	ID          int64  `json:"id"`
	TramiteID   int64  `json:"tramite_id"`
	Alias       string `json:"alias"`
	Tipo        string `json:"tipo"`
	Normalizado string `json:"normalizado"`
}

type TramiteVersion struct {
	ID                    int64      `json:"id"`
	TramiteID             int64      `json:"tramite_id"`
	NumeroVersion         int        `json:"numero_version"`
	EstadoEditorial       string     `json:"estado_editorial"`
	Titulo                string     `json:"titulo"`
	Resumen               string     `json:"resumen"`
	Descripcion           string     `json:"descripcion"`
	PublicoObjetivo       *string    `json:"publico_objetivo,omitempty"`
	ResultadoEsperado     *string    `json:"resultado_esperado,omitempty"`
	Advertencias          *string    `json:"advertencias,omitempty"`
	PlazoTexto            *string    `json:"plazo_texto,omitempty"`
	VigenciaResultadoTexto *string   `json:"vigencia_resultado_texto,omitempty"`
	RequiereCita          bool       `json:"requiere_cita"`
	URLInicio             *string    `json:"url_inicio,omitempty"`
	NotasInternas         *string    `json:"notas_internas,omitempty"`
	ValidoDesde           *time.Time `json:"valido_desde,omitempty"`
	ValidoHasta           *time.Time `json:"valido_hasta,omitempty"`
	VerificadoEn          *time.Time `json:"verificado_en,omitempty"`
	ProximaRevisionEn     *time.Time `json:"proxima_revision_en,omitempty"`
	CreadoPor             *int64     `json:"creado_por,omitempty"`
	RevisadoPor           *int64     `json:"revisado_por,omitempty"`
	AprobadoPor           *int64     `json:"aprobado_por,omitempty"`
	CreadoEn              time.Time  `json:"creado_en"`
	PublicadoEn           *time.Time `json:"publicado_en,omitempty"`
	HashContenido         *string    `json:"hash_contenido,omitempty"`
}

// TramitePublicado es la vista vw_tramite_publicado aplanada para lectura pública
type TramitePublicado struct {
	ID                    int64      `json:"id"`
	Codigo                string     `json:"codigo"`
	Slug                  string     `json:"slug"`
	InstitucionID         int64      `json:"institucion_id"`
	Institucion           string     `json:"institucion"`
	Sigla                 string     `json:"sigla"`
	CategoriaID           *int64     `json:"categoria_id,omitempty"`
	Categoria             *string    `json:"categoria,omitempty"`
	VersionID             int64      `json:"version_id"`
	NumeroVersion         int        `json:"numero_version"`
	Titulo                string     `json:"titulo"`
	Resumen               string     `json:"resumen"`
	Descripcion           string     `json:"descripcion"`
	PublicoObjetivo       *string    `json:"publico_objetivo,omitempty"`
	ResultadoEsperado     *string    `json:"resultado_esperado,omitempty"`
	Advertencias          *string    `json:"advertencias,omitempty"`
	PlazoTexto            *string    `json:"plazo_texto,omitempty"`
	VigenciaResultadoTexto *string   `json:"vigencia_resultado_texto,omitempty"`
	RequiereCita          bool       `json:"requiere_cita"`
	URLInicio             *string    `json:"url_inicio,omitempty"`
	VerificadoEn          *time.Time `json:"verificado_en,omitempty"`
	ProximaRevisionEn     *time.Time `json:"proxima_revision_en,omitempty"`
	PublicadoEn           *time.Time `json:"publicado_en,omitempty"`
}

type FuenteOficial struct {
	ID              int64      `json:"id"`
	InstitucionID   *int64     `json:"institucion_id,omitempty"`
	Tipo            string     `json:"tipo"`
	Titulo          string     `json:"titulo"`
	URL             string     `json:"url"`
	FechaDocumento  *string    `json:"fecha_documento,omitempty"`
	AutoridadEmisora *string   `json:"autoridad_emisora,omitempty"`
	HashContenido   *string    `json:"hash_contenido,omitempty"`
	ConsultadoEn    time.Time  `json:"consultado_en"`
	Estado          string     `json:"estado"`
	EsOficial       bool       `json:"es_oficial"`
}

type ModalidadTramite struct {
	ID               int64   `json:"id"`
	TramiteVersionID int64   `json:"tramite_version_id"`
	Tipo             string  `json:"tipo"`
	Nombre           string  `json:"nombre"`
	Descripcion      *string `json:"descripcion,omitempty"`
	URLInicio        *string `json:"url_inicio,omitempty"`
	RequiereCita     bool    `json:"requiere_cita"`
	Orden            int     `json:"orden"`
}

type RequisitoTramite struct {
	ID                int64           `json:"id"`
	TramiteVersionID  int64           `json:"tramite_version_id"`
	ModalidadID       *int64          `json:"modalidad_id,omitempty"`
	Codigo            string          `json:"codigo"`
	Titulo            string          `json:"titulo"`
	Descripcion       *string         `json:"descripcion,omitempty"`
	Obligatorio       bool            `json:"obligatorio"`
	CantidadOriginales int            `json:"cantidad_originales"`
	CantidadCopias    int            `json:"cantidad_copias"`
	Formato           *string         `json:"formato,omitempty"`
	VigenciaDocumento *string         `json:"vigencia_documento,omitempty"`
	Emisor            *string         `json:"emisor,omitempty"`
	GrupoAlternativa  *string         `json:"grupo_alternativa,omitempty"`
	AplicaSi          json.RawMessage `json:"aplica_si"`
	Orden             int             `json:"orden"`
}

type PasoTramite struct {
	ID               int64   `json:"id"`
	TramiteVersionID int64   `json:"tramite_version_id"`
	ModalidadID      *int64  `json:"modalidad_id,omitempty"`
	Numero           int     `json:"numero"`
	Titulo           string  `json:"titulo"`
	Descripcion      string  `json:"descripcion"`
	Lugar            *string `json:"lugar,omitempty"`
	URLAccion        *string `json:"url_accion,omitempty"`
	EsFueraSistema   bool    `json:"es_fuera_sistema"`
}

type CostoTramite struct {
	ID               int64           `json:"id"`
	TramiteVersionID int64           `json:"tramite_version_id"`
	ModalidadID      *int64          `json:"modalidad_id,omitempty"`
	Concepto         string          `json:"concepto"`
	Moneda           string          `json:"moneda"`
	Monto            *float64        `json:"monto,omitempty"`
	MontoDesde       *float64        `json:"monto_desde,omitempty"`
	MontoHasta       *float64        `json:"monto_hasta,omitempty"`
	EsGratuito       bool            `json:"es_gratuito"`
	Periodicidad     *string         `json:"periodicidad,omitempty"`
	MedioPago        *string         `json:"medio_pago,omitempty"`
	Instrucciones    *string         `json:"instrucciones,omitempty"`
	AplicaSi         json.RawMessage `json:"aplica_si"`
	VerificadoEn     *time.Time      `json:"verificado_en,omitempty"`
}

type ResultadoTramite struct {
	ID               int64   `json:"id"`
	TramiteVersionID int64   `json:"tramite_version_id"`
	Nombre           string  `json:"nombre"`
	Formato          *string `json:"formato,omitempty"`
	Vigencia         *string `json:"vigencia,omitempty"`
	Entrega          *string `json:"entrega,omitempty"`
	Orden            int     `json:"orden"`
}

// ─── Usuarios y sesiones admin ────────────────────────────────────────────────

type UsuarioAdmin struct {
	ID                  int64      `json:"id"`
	Correo              string     `json:"correo"`
	Nombre              string     `json:"nombre"`
	PasswordHash        *string    `json:"-"` // nunca serializar
	ProveedorIdentidad  string     `json:"proveedor_identidad"`
	SujetoExterno       *string    `json:"sujeto_externo,omitempty"`
	Estado              string     `json:"estado"`
	UltimoAccesoEn      *time.Time `json:"ultimo_acceso_en,omitempty"`
	CreadoEn            time.Time  `json:"creado_en"`
	ActualizadoEn       time.Time  `json:"actualizado_en"`
}

type SesionAdmin struct {
	ID               string     `json:"id"`
	UsuarioID        int64      `json:"usuario_id"`
	RefreshTokenHash string     `json:"-"`
	IPHash           *string    `json:"-"`
	UserAgent        *string    `json:"user_agent,omitempty"`
	ExpiraEn         time.Time  `json:"expira_en"`
	RevocadaEn       *time.Time `json:"revocada_en,omitempty"`
	CreadaEn         time.Time  `json:"creada_en"`
	UltimoUsoEn      time.Time  `json:"ultimo_uso_en"`
}

type Rol struct {
	ID       int             `json:"id"`
	Codigo   string          `json:"codigo"`
	Nombre   string          `json:"nombre"`
	Permisos json.RawMessage `json:"permisos"`
}

type UsuarioRol struct {
	ID            int64      `json:"id"`
	UsuarioID     int64      `json:"usuario_id"`
	RolID         int        `json:"rol_id"`
	InstitucionID *int64     `json:"institucion_id,omitempty"`
	AsignadoEn    time.Time  `json:"asignado_en"`
}

// ─── Auditoría ─────────────────────────────────────────────────────────────────

type EventoAuditoria struct {
	ID           int64           `json:"id"`
	ActorID      *int64          `json:"actor_id,omitempty"`
	Accion       string          `json:"accion"`
	EntidadTipo  string          `json:"entidad_tipo"`
	EntidadID    *string         `json:"entidad_id,omitempty"`
	Antes        json.RawMessage `json:"antes,omitempty"`
	Despues      json.RawMessage `json:"despues,omitempty"`
	IPHash       *string         `json:"-"`
	UserAgent    *string         `json:"user_agent,omitempty"`
	OcurridoEn   time.Time       `json:"ocurrido_en"`
}

// ─── Ingesta ───────────────────────────────────────────────────────────────────

type FuenteIngesta struct {
	ID                  int64           `json:"id"`
	InstitucionID       *int64          `json:"institucion_id,omitempty"`
	Nombre              string          `json:"nombre"`
	Tipo                string          `json:"tipo"`
	URL                 string          `json:"url"`
	FrecuenciaCron      *string         `json:"frecuencia_cron,omitempty"`
	Configuracion       json.RawMessage `json:"configuracion"`
	Estado              string          `json:"estado"`
	UltimaEjecucionEn   *time.Time      `json:"ultima_ejecucion_en,omitempty"`
	ProximaEjecucionEn  *time.Time      `json:"proxima_ejecucion_en,omitempty"`
	CreadoEn            time.Time       `json:"creado_en"`
}

type EjecucionIngesta struct {
	ID               string     `json:"id"`
	FuenteIngestaID  int64      `json:"fuente_ingesta_id"`
	Estado           string     `json:"estado"`
	HTTPStatus       *int       `json:"http_status,omitempty"`
	HashContenido    *string    `json:"hash_contenido,omitempty"`
	RegistrosLeidos  int        `json:"registros_leidos"`
	CandidatosCreados int       `json:"candidatos_creados"`
	MensajeError     *string    `json:"mensaje_error,omitempty"`
	IniciadaEn       time.Time  `json:"iniciada_en"`
	FinalizadaEn     *time.Time `json:"finalizada_en,omitempty"`
}

type CandidatoIngesta struct {
	ID                string          `json:"id"`
	EjecucionID       string          `json:"ejecucion_id"`
	TramiteIDSugerido *int64          `json:"tramite_id_sugerido,omitempty"`
	DatosExtraidos    json.RawMessage `json:"datos_extraidos"`
	TextoCrudo        *string         `json:"texto_crudo,omitempty"`
	HashContenido     string          `json:"hash_contenido"`
	Confianza         *float64        `json:"confianza,omitempty"`
	Estado            string          `json:"estado"`
	RevisadoPor       *int64          `json:"revisado_por,omitempty"`
	RevisadoEn        *time.Time      `json:"revisado_en,omitempty"`
}

// ─── RAG y conversaciones ─────────────────────────────────────────────────────

type FragmentoConocimiento struct {
	ID               string          `json:"id"`
	TramiteVersionID int64           `json:"tramite_version_id"`
	Tipo             string          `json:"tipo"`
	ReferenciaID     *int64          `json:"referencia_id,omitempty"`
	Contenido        string          `json:"contenido"`
	Metadatos        json.RawMessage `json:"metadatos"`
	HashContenido    string          `json:"hash_contenido"`
	ModeloEmbedding  *string         `json:"modelo_embedding,omitempty"`
	EstadoEmbedding  string          `json:"estado_embedding"`
	GeneradoEn       *time.Time      `json:"generado_en,omitempty"`
}

type TrabajoEmbedding struct {
	ID              string     `json:"id"`
	FragmentoID     string     `json:"fragmento_id"`
	Modelo          string     `json:"modelo"`
	Estado          string     `json:"estado"`
	Intentos        int        `json:"intentos"`
	ErrorDetalle    *string    `json:"error_detalle,omitempty"`
	DisponibleDesde time.Time  `json:"disponible_desde"`
	IniciadoEn      *time.Time `json:"iniciado_en,omitempty"`
	CompletadoEn    *time.Time `json:"completado_en,omitempty"`
}

type SesionAnonima struct {
	ID                     string     `json:"id"`
	IdentificadorHash      string     `json:"-"`
	Idioma                 string     `json:"idioma"`
	MunicipioID            *int64     `json:"municipio_id,omitempty"`
	ConsentimientoAnalytics bool      `json:"consentimiento_analytics"`
	CreadaEn               time.Time  `json:"creada_en"`
	ExpiraEn               time.Time  `json:"expira_en"`
	UltimaActividadEn      time.Time  `json:"ultima_actividad_en"`
}

type Conversacion struct {
	ID         string     `json:"id"`
	SesionID   string     `json:"sesion_id"`
	Titulo     *string    `json:"titulo,omitempty"`
	Estado     string     `json:"estado"`
	CreadaEn   time.Time  `json:"creada_en"`
	CerradaEn  *time.Time `json:"cerrada_en,omitempty"`
}

type MensajeConversacion struct {
	ID             string    `json:"id"`
	ConversacionID string    `json:"conversacion_id"`
	Rol            string    `json:"rol"`
	Contenido      string    `json:"contenido"`
	Modelo         *string   `json:"modelo,omitempty"`
	TokensEntrada  *int      `json:"tokens_entrada,omitempty"`
	TokensSalida   *int      `json:"tokens_salida,omitempty"`
	LatenciaMs     *int      `json:"latencia_ms,omitempty"`
	EstadoSeguridad string   `json:"estado_seguridad"`
	CreadoEn       time.Time `json:"creado_en"`
}

type CitaRAG struct {
	RecuperacionID    string   `json:"recuperacion_id"`
	FragmentoID       string   `json:"fragmento_id"`
	Posicion          int      `json:"posicion"`
	Similitud         *float64 `json:"similitud,omitempty"`
	IncluidaEnPrompt  bool     `json:"incluida_en_prompt"`
	// Campos enriquecidos para respuesta al cliente
	TramiteSlug       *string  `json:"tramite_slug,omitempty"`
	TramiteTitulo     *string  `json:"tramite_titulo,omitempty"`
	TipoFragmento     *string  `json:"tipo_fragmento,omitempty"`
	URLFuente         *string  `json:"url_fuente,omitempty"`
}

// ─── Reportes ciudadanos ───────────────────────────────────────────────────────

type ReporteCiudadano struct {
	ID                    string     `json:"id"`
	TramiteID             *int64     `json:"tramite_id,omitempty"`
	OficinaID             *int64     `json:"oficina_id,omitempty"`
	Tipo                  string     `json:"tipo"`
	Descripcion           string     `json:"descripcion"`
	CorreoContactoCifrado []byte     `json:"-"`
	Estado                string     `json:"estado"`
	AsignadoA             *int64     `json:"asignado_a,omitempty"`
	CreadoEn              time.Time  `json:"creado_en"`
	ResueltoEn            *time.Time `json:"resuelto_en,omitempty"`
}

// ─── DTOs de respuesta de la API pública ──────────────────────────────────────

// TramiteListItem es la representación reducida en listados
type TramiteListItem struct {
	ID                int64      `json:"id"`
	Slug              string     `json:"slug"`
	Titulo            string     `json:"titulo"`
	Resumen           string     `json:"resumen"`
	Institucion       string     `json:"institucion"`
	InstitucionSigla  string     `json:"institucion_sigla"`
	Categoria         *string    `json:"categoria,omitempty"`
	VerificadoEn      *time.Time `json:"verificado_en,omitempty"`
	ProximaRevisionEn *time.Time `json:"proxima_revision_en,omitempty"`
	RequiereVerificacion bool    `json:"requiere_verificacion"`
}

// TramiteFicha es la respuesta completa de un trámite individual
type TramiteFicha struct {
	ID                    int64              `json:"id"`
	Slug                  string             `json:"slug"`
	Titulo                string             `json:"titulo"`
	Resumen               string             `json:"resumen"`
	Descripcion           string             `json:"descripcion"`
	PublicoObjetivo       *string            `json:"publico_objetivo,omitempty"`
	ResultadoEsperado     *string            `json:"resultado_esperado,omitempty"`
	Advertencias          *string            `json:"advertencias,omitempty"`
	PlazoTexto            *string            `json:"plazo_texto,omitempty"`
	VigenciaResultadoTexto *string           `json:"vigencia_resultado_texto,omitempty"`
	RequiereCita          bool               `json:"requiere_cita"`
	URLInicio             *string            `json:"url_inicio,omitempty"`
	VerificadoEn          *time.Time         `json:"verificado_en,omitempty"`
	ProximaRevisionEn     *time.Time         `json:"proxima_revision_en,omitempty"`
	RequiereVerificacion  bool               `json:"requiere_verificacion"`
	Institucion           InstitucionResumen `json:"institucion"`
	Categoria             *string            `json:"categoria,omitempty"`
	Modalidades           []ModalidadDetalle `json:"modalidades"`
	Requisitos            []RequisitoTramite `json:"requisitos"`
	Pasos                 []PasoTramite      `json:"pasos"`
	Costos                []CostoTramite     `json:"costos"`
	Resultados            []ResultadoTramite `json:"resultados"`
	Fuentes               []FuenteResumen    `json:"fuentes"`
}

type InstitucionResumen struct {
	Codigo   string  `json:"codigo"`
	Nombre   string  `json:"nombre"`
	Sigla    string  `json:"sigla"`
	SitioWeb *string `json:"sitio_web,omitempty"`
}

type ModalidadDetalle struct {
	ID           int64    `json:"id"`
	Tipo         string   `json:"tipo"`
	Nombre       string   `json:"nombre"`
	Descripcion  *string  `json:"descripcion,omitempty"`
	URLInicio    *string  `json:"url_inicio,omitempty"`
	RequiereCita bool     `json:"requiere_cita"`
	Orden        int      `json:"orden"`
	Oficinas     []OficinaSummary `json:"oficinas,omitempty"`
}

type FuenteResumen struct {
	Tipo   string  `json:"tipo"`
	Titulo string  `json:"titulo"`
	URL    string  `json:"url"`
	Uso    string  `json:"uso"`
}

// ─── JWT Claims ────────────────────────────────────────────────────────────────

type AdminClaims struct {
	UsuarioID int64    `json:"usuario_id"`
	Correo    string   `json:"correo"`
	Permisos  []string `json:"permisos"`
}
