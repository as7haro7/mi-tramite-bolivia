# Guion visual ampliado — Mi Trámite Bolivia

Fuente principal: `Proyecto_Final_Mi_Tramite_Bolivia_INF264.docx`.

> Nota: no se adjuntó archivo Figma ni capturas del prototipo. La diapositiva 15 usa maquetas reconstruidas desde la Tabla 8 y lo declara expresamente.

## Diapositiva 1

**Título de la diapositiva**
Mi Trámite Bolivia

**Textos clave o viñetas**
• Tu asistente inteligente para realizar trámites públicos.
• Proyecto Final INF-264 — UMSA, Carrera de Informática.
• Equipo: Canqui Phuña Helen Yvette Cecilia; Chana Saico Cecilia; Enriquez Aduviri Vanesa Alejandra; Mendoza Mamani Ricardo Einer; Erick Fernando Poma Condori.

**Indicaciones visuales precisas**
Portada azul marino con logotipo del proyecto como foco y escudo UMSA en el panel lateral. Los nombres deben quedar completos, sin abreviaturas.

**Instrucciones específicas para diagramas/tablas**
No requiere diagrama. Mantener marca y promesa a la izquierda; datos académicos y equipo a la derecha.

## Diapositiva 2

**Título de la diapositiva**
Árbol del problema

**Textos clave o viñetas**
• Causas: información dispersa, lenguaje técnico y variación según institución, trámite y caso.
• Problema central: baja capacidad de preparación correcta antes del primer contacto institucional.
• Efectos: documentación incompleta, costos evitables e incertidumbre.

**Indicaciones visuales precisas**
Diagrama causal puro, sin párrafos largos. Causas arriba en turquesa; problema central en azul marino; efectos abajo en dorado.

**Instrucciones específicas para diagramas/tablas**
Tres tarjetas superiores conectan al problema central; este conecta con tres efectos. La pregunta de investigación ocupa una banda inferior.

## Diapositiva 3

**Título de la diapositiva**
Objetivos y métricas de validación

**Textos clave o viñetas**
• Objetivo general: diseñar y desarrollar un MVP móvil inteligente para orientación sobre trámites públicos bolivianos.
• Objetivos específicos resumidos: investigación, catálogo versionado, aplicación Flutter, API Go, IA intercambiable y validación.
• Indicadores: preparación completa, resolución en primer intento, tiempo de orientación, exactitud, respuestas con fuente y satisfacción.

**Indicaciones visuales precisas**
Objetivo general en una banda superior; cinco objetivos específicos como tarjetas; seis indicadores en una matriz 3×2.

**Instrucciones específicas para diagramas/tablas**
No inventar metas porcentuales. Los indicadores se presentan como variables que medirán el piloto.

## Diapositiva 4

**Título de la diapositiva**
Descripción del emprendimiento

**Textos clave o viñetas**
• La plataforma traduce información institucional a guías verificables y personalizadas.
• Cada guía contiene pasos, requisitos, costos, modalidad, horarios, ubicación, fuentes, fecha y advertencias.
• La aplicación orienta y deriva; no sustituye ni representa a la institución.
• El MVP comienza con 20–30 trámites de alta demanda.

**Indicaciones visuales precisas**
Flujo horizontal de cinco etapas seguido por tres bloques de capacidades.

**Instrucciones específicas para diagramas/tablas**
Secuencia: Describir → Aclarar → Orientar → Preparar → Derivar. La derivación se diferencia en dorado.

## Diapositiva 5

**Título de la diapositiva**
Requisitos del sistema — Funcionales

**Textos clave o viñetas**
• RF-01: Buscar trámites por texto, categoría o institución. Prioridad alta.
• RF-02: Visualizar ficha con fuente y fecha de última verificación. Prioridad alta.
• RF-03: Conversar con el asistente y obtener una guía personalizada. Prioridad alta.
• RF-04: Generar y guardar checklist de documentos. Prioridad alta.
• RF-05: Registrar reportes de información posiblemente desactualizada. Prioridad media.
• RF-06: Administrar versiones y estados de publicación de trámites. Prioridad alta.

**Indicaciones visuales precisas**
Tabla limpia con seis filas; alternar fondos y destacar evidencia, asistente y versionado.

**Instrucciones específicas para diagramas/tablas**
Columnas: Código | Requisito funcional | Prioridad. Mantener los textos exactos de la Tabla 3.

## Diapositiva 6

**Título de la diapositiva**
Requisitos del sistema — No funcionales

**Textos clave o viñetas**
• RNF-01: respuesta de API inferior a 2 segundos en consultas no generativas.
• RNF-02: disponibilidad objetivo de 99 % durante el piloto.
• RNF-03: cifrado en tránsito, secretos, RBAC y auditoría.
• RNF-04: interfaz accesible y utilizable con conexión limitada.
• RNF-05: portabilidad entre proveedores de IA y entornos de despliegue.

**Indicaciones visuales precisas**
Cinco tarjetas métricas con una tarjeta final que explica la decisión de portabilidad.

**Instrucciones específicas para diagramas/tablas**
Dar máxima jerarquía visual a “< 2 s” y “PORTABLE”, porque son requisitos explícitamente solicitados para la defensa.

## Diapositiva 7

**Título de la diapositiva**
Justificación

**Textos clave o viñetas**
• Técnica: Flutter, Go, PostgreSQL, Render, Neon y APIs de IA son suficientes para un MVP.
• Económica: una base compartida y servicios administrados reducen duplicación.
• Social: disminuye asimetrías de información y errores evitables.
• Académica: integra emprendimiento, UX, software, datos, IA, seguridad y economía.

**Indicaciones visuales precisas**
Cuatro tarjetas grandes con iconos numéricos y una ecuación de innovación al cierre.

**Instrucciones específicas para diagramas/tablas**
La tarjeta final debe remarcar que el modelo de lenguaje no es la fuente.

## Diapositiva 8

**Título de la diapositiva**
Público objetivo

**Textos clave o viñetas**
• Camila: estudiante que necesita fuente vigente y checklist.
• José: trabajador independiente que planifica costos, tiempos y secuencias.
• María: cuidadora familiar con conectividad intermitente.
• Segmentos: ciudadanos, emprendedores, estudiantes, familias e instituciones.
• Early adopters: comunidad UMSA y emprendedores de La Paz.

**Indicaciones visuales precisas**
Tres tarjetas de usuario con iconos de inicial grandes y cinco tarjetas de segmentos.

**Instrucciones específicas para diagramas/tablas**
Resaltar instituciones en dorado para mostrar la coexistencia B2C/B2B.

## Diapositiva 9

**Título de la diapositiva**
Modelo de negocio — Canvas

**Textos clave o viñetas**
• Modelo freemium con información básica gratuita.
• Premium: recordatorios, perfiles familiares, historial y alertas.
• Instituciones: panel, estadísticas, configuración y white-label.
• No vender datos personales ni cobrar por información pública.

**Indicaciones visuales precisas**
Cuadrícula visual de nueve bloques según el Business Model Canvas clásico.

**Instrucciones específicas para diagramas/tablas**
Orden: Socios | Actividades/Recursos | Propuesta | Relación/Canales | Segmentos; abajo Costos e Ingresos.

## Diapositiva 10

**Título de la diapositiva**
Stack tecnológico del MVP

**Textos clave o viñetas**
• Frontend: Flutter/Dart, Riverpod, Dio y Drift.
• Backend: Go/Gin, pgx, sqlc y Goose.
• Datos: PostgreSQL en Neon, full-text y pgvector opcional.
• IA: Gemini y Qwen mediante RAG y JSON Schema.
• Operación: Render, Docker, GitHub Actions y OpenAPI.

**Indicaciones visuales precisas**
Cinco columnas de stack y cinco principios arquitectónicos debajo.

**Instrucciones específicas para diagramas/tablas**
Esta diapositiva enumera tecnologías; la arquitectura lógica se muestra por separado.

## Diapositiva 11

**Título de la diapositiva**
Arquitectura lógica propuesta del MVP

**Textos clave o viñetas**
• Clientes: aplicación Flutter y panel Flutter Web.
• Comunicación REST/HTTPS con un backend Go.
• Backend: endpoints, monolito modular y adaptadores.
• Infraestructura: Neon, proveedores Gemini/Qwen y Render/CI/CD.
• La aplicación nunca se conecta directamente a la base de datos ni a la IA.

**Indicaciones visuales precisas**
Diagrama exclusivo de arquitectura, separado de la lista de tecnologías.

**Instrucciones específicas para diagramas/tablas**
Izquierda: clientes y backend por capas. Derecha: tres servicios externos. Conectar únicamente desde los adaptadores.

## Diapositiva 12

**Título de la diapositiva**
Modelo de datos y seguridad

**Textos clave o viñetas**
• Entidades: institution, procedure, procedure_version, requirement, requirement_rule, source, location, user_checklist, conversation, feedback y audit_log.
• Las versiones publicadas son inmutables y conservan historial.
• Seguridad: OWASP API Top 10, RBAC, JWT, 2FA administrativa, TLS, secretos, límites, auditoría y backups.
• Privacidad: minimización, consultas sin cuenta, sin imágenes de documentos y eliminación de historial.

**Indicaciones visuales precisas**
Agrupar las entidades en cuatro dominios y colocar un panel vertical de seguridad.

**Instrucciones específicas para diagramas/tablas**
El panel de seguridad debe ser visualmente dominante y mostrar OWASP, RBAC y JWT.

## Diapositiva 13

**Título de la diapositiva**
El corazón de la IA: RAG controlado

**Textos clave o viñetas**
• Flujo: consulta → recuperación del catálogo → generación → control → guía.
• Gemini y Qwen se integran detrás de una interfaz común.
• JSON Schema obliga a producir campos estructurados.
• Abstención: si las fuentes no permiten confirmar, la IA no inventa.
• Cada requisito debe poder trazarse al contexto recuperado.

**Indicaciones visuales precisas**
Flujo horizontal de cinco pasos sobre fondo azul marino; dos tarjetas grandes para JSON Schema y abstención.

**Instrucciones específicas para diagramas/tablas**
La salida final debe decir explícitamente “guía o abstención”, no solo “respuesta”.

## Diapositiva 14

**Título de la diapositiva**
Visión futura — Migración a un LLM local

**Textos clave o viñetas**
• Plan: ejecutar un modelo Qwen instructivo mediante vLLM en infraestructura controlada.
• Objetivos: privacidad, control y reducción potencial de costos a escala.
• Antes de migrar se medirán costo, calidad en español boliviano, errores, privacidad, latencia y demanda.
• Se requieren GPU, monitoreo, colas, concurrencia, continuidad y respaldo externo.

**Indicaciones visuales precisas**
Comparativa visual entre etapa API y visión local; debajo, cinco criterios y una banda de requisitos.

**Instrucciones específicas para diagramas/tablas**
No presentar la migración como decisión inmediata: debe aparecer condicionada a métricas.

## Diapositiva 15

**Título de la diapositiva**
Prototipo y flujo de experiencia

**Textos clave o viñetas**
• Flujo: Inicio → Consultar/Buscar → Revisar requisitos → Preparar checklist → Derivar al canal oficial.
• Pantallas destacadas: Inicio, Asistente y Checklist.
• La ficha y las respuestas muestran fuente y fecha.
• La aplicación no afirma que completó el trámite.

**Indicaciones visuales precisas**
Flujo de la Figura 3 arriba y tres mockups de teléfono debajo.

**Instrucciones específicas para diagramas/tablas**
No se encontró un archivo Figma ni capturas reales. Las maquetas se reconstruyen desde la Tabla 8 y se etiquetan como tales.

## Diapositiva 16

**Título de la diapositiva**
Análisis de mercado y competencia

**Textos clave o viñetas**
• Alternativas: gob.bo, Ciudadanía Digital/PTC, portales, buscadores/redes y tramitadores.
• Mi Trámite se posiciona como complemento interoperable.
• Diferenciales: personalización, trazabilidad, móvil, IA controlada y enfoque nacional.

**Indicaciones visuales precisas**
Tabla de seis filas con Mi Trámite Bolivia resaltado.

**Instrucciones específicas para diagramas/tablas**
Columnas: Alternativa | Fortaleza | Brecha/posición. Cerrar con una banda de diferenciadores.

## Diapositiva 17

**Título de la diapositiva**
Estrategia de captación

**Textos clave o viñetas**
• Alianzas con centros de estudiantes UMSA, incubadoras, ferias y cámaras.
• Contenido educativo en TikTok, Facebook e Instagram.
• Códigos QR en materiales de aliados, con autorización.
• Programa beta con funciones premium temporales.
• Páginas públicas indexables por trámite.

**Indicaciones visuales precisas**
Cinco tarjetas de canales y una secuencia de entrada al mercado.

**Instrucciones específicas para diagramas/tablas**
Ruta: La Paz → 20–30 trámites → piloto UMSA → demo B2B → expansión por paquetes.

## Diapositiva 18

**Título de la diapositiva**
Estudio económico — Inversión y operación

**Textos clave o viñetas**
• Valor económico total: 97.350 Bs.
• Desembolso estimado con desarrollo fundador: 24.500 Bs.
• Costo sostenible: 6.000 Bs./mes.
• En fase académica, aportes de curación y soporte pueden reducir el efectivo a aproximadamente 2.500 Bs./mes.

**Indicaciones visuales precisas**
Tres cifras principales; barras de inversión a la izquierda; costo mensual apilado al centro; lectura estratégica a la derecha.

**Instrucciones específicas para diagramas/tablas**
La inversión debe diferenciar valor económico azul y desembolso dorado.

## Diapositiva 19

**Título de la diapositiva**
Estudio económico — Ingresos y equilibrio

**Textos clave o viñetas**
• Premium: 20 Bs./mes.
• Licencia institucional: 2.000 Bs./mes.
• White-label: 8.000–20.000 Bs.
• Equilibrio: 4 instituciones; o 2 instituciones + 134 premium; o 334 premium.
• Escenario base: 3 instituciones y 250 premium hacia el último trimestre.

**Indicaciones visuales precisas**
Cuatro tarjetas de ingresos y tres filas grandes de punto de equilibrio.

**Instrucciones específicas para diagramas/tablas**
Usar exactamente las contribuciones de la Tabla 13: 1.800 Bs. por licencia y 18 Bs. por premium.

## Diapositiva 20

**Título de la diapositiva**
Mapa de calor de riesgos

**Textos clave o viñetas**
• Alta/alta: información desactualizada.
• Alta/media: alcance excesivo.
• Media/alta: alucinación, acceso, confusión oficial, baja adopción y falta de actualización.
• Media/media: caída de IA, costos variables y brecha digital.

**Indicaciones visuales precisas**
Mapa de calor 3×3 verde, amarillo y rojo con los diez riesgos dentro de sus celdas.

**Instrucciones específicas para diagramas/tablas**
Eje vertical: probabilidad baja/media/alta. Eje horizontal: impacto bajo/medio/alto. Panel lateral de mitigaciones.

## Diapositiva 21

**Título de la diapositiva**
Cronograma de implementación

**Textos clave o viñetas**
• Duración total: 20 semanas.
• Sprints de dos semanas con planificación, desarrollo, pruebas, demo y retrospectiva.
• El catálogo avanza en paralelo al software.

**Indicaciones visuales precisas**
Gantt del documento a pantalla completa con una banda de hitos.

**Instrucciones específicas para diagramas/tablas**
Respetar exactamente las nueve fases y sus semanas.

## Diapositiva 22

**Título de la diapositiva**
Entregables e hitos

**Textos clave o viñetas**
• H1 Validación: problema y público confirmados.
• H2 Prototipo: flujo navegable probado.
• H3 Alfa: búsqueda, ficha y checklist.
• H4 Asistente: fuentes y evaluación.
• H5 Beta: pruebas críticas y monitoreo.
• H6 Piloto: métricas recopiladas.
• H7 Lanzamiento/defensa: MVP, documentación y continuidad.

**Indicaciones visuales precisas**
Siete tarjetas de hitos con semana, definición y entregable.

**Instrucciones específicas para diagramas/tablas**
Dar mayor énfasis visual a Alfa, Beta, Piloto y Lanzamiento mediante colores distintos.

## Diapositiva 23

**Título de la diapositiva**
Conclusiones y recomendaciones

**Textos clave o viñetas**
• El proyecto es técnicamente alcanzable y socialmente relevante.
• El catálogo confiable es el principal activo.
• La viabilidad empresarial requiere ingresos institucionales.
• Recomendaciones: limitar el piloto, gobernar contenido, validar B2B y medir antes de migrar a Qwen local.

**Indicaciones visuales precisas**
Cierre con una afirmación, tres conclusiones y cuatro recomendaciones accionables.

**Instrucciones específicas para diagramas/tablas**
Terminar con la ecuación: menos incertidumbre + preparación correcta + evidencia verificable.
