# Guion visual — Mi Trámite Bolivia

Fuente exclusiva: `Proyecto_Final_Mi_Tramite_Bolivia_INF264.docx`.

## 1. Mi Trámite Bolivia

**Textos clave o viñetas**

- “Tu asistente inteligente para realizar trámites públicos”.
- Proyecto Final de INF-264 — Emprendimiento e Innovación Tecnológica.
- Universidad Mayor de San Andrés, Carrera de Informática; La Paz, Bolivia — 2026.

**Indicaciones visuales precisas**

Portada sobria en azul marino. El logotipo incluido en el documento es el foco principal y el escudo UMSA ocupa un panel lateral. La franja turquesa y el acento dorado refuerzan confianza, tecnología y ciudadanía.

**Instrucciones específicas para diagramas/tablas**

No requiere diagrama. Jerarquía: marca y promesa al centro-izquierda; datos académicos y equipo en el panel derecho.

## 2. Identificación del problema

**Textos clave o viñetas**

- El ciudadano debe identificar institución, trámite, requisitos aplicables, costos, horarios, modalidad y respuesta ante observaciones.
- La información dispersa y el lenguaje técnico generan documentación incompleta, filas, visitas repetidas y costos.
- Problema central: baja capacidad de preparación correcta antes del primer contacto con la institución.
- Hipótesis: una guía personalizada con checklist, costo, ubicación, fuente y vigencia mejora la preparación y la resolución en el primer intento.

**Indicaciones visuales precisas**

Árbol causal minimalista: tres causas arriba, problema central en el centro y tres efectos abajo. Turquesa para causas, azul marino para el núcleo y dorado para efectos.

**Instrucciones específicas para diagramas/tablas**

Conectar las causas “información dispersa”, “lenguaje técnico” y “ruta variable” con el problema central. Desde este bloque, conectar “costo directo”, “tiempo perdido” y “costo emocional”. En una banda inferior, mostrar la hipótesis y los seis indicadores: preparación completa, primer intento, tiempo de orientación, exactitud, respuestas con fuente y satisfacción.

## 3. Descripción del emprendimiento

**Textos clave o viñetas**

- Promesa: antes de hacer fila, el ciudadano sabrá qué preparar, dónde ir y qué verificar.
- La guía incluye requisitos obligatorios y condicionales, pasos, costo, modalidad, horarios, ubicación, fuentes, fecha y advertencias.
- El asistente se abstiene cuando no puede confirmar una respuesta y deriva al contacto oficial.
- El MVP se limita a 20–30 trámites de alta demanda.

**Indicaciones visuales precisas**

Flujo horizontal de cinco etapas y tres bloques de capacidades: MVP ciudadano, arquitectura de confianza y límite responsable.

**Instrucciones específicas para diagramas/tablas**

Orden: Describir necesidad → Aclarar caso → Orientar con tarjetas → Preparar checklist → Derivar al canal oficial. La última etapa debe distinguirse en dorado para remarcar que la aplicación no ejecuta el trámite.

## 4. Justificación

**Textos clave o viñetas**

- Técnica: Flutter, Go, PostgreSQL, Render, Neon y APIs de IA permiten implementar el MVP con componentes maduros.
- Económica: una base multiplataforma y servicios administrados reducen duplicación y permiten comenzar con costos variables.
- Social: puede reducir asimetrías de información, errores y desplazamientos evitables.
- Académica e innovadora: integra emprendimiento, UX, ingeniería, datos, IA, seguridad y evaluación económica.
- El modelo no es la fuente: consulta un catálogo oficial versionado.

**Indicaciones visuales precisas**

Cuatro cuadrantes alrededor de la afirmación “reducir incertidumbre con evidencia”. Cada cuadrante usa un color de acento distinto.

**Instrucciones específicas para diagramas/tablas**

Técnica y Económica arriba; Social y Académica/Innovadora abajo. Conectar las tarjetas al centro. Añadir el sello “fuente + vigencia + abstención”.

## 5. Público objetivo

**Textos clave o viñetas**

- Segmentos: ciudadanos de 18–45 años, emprendedores, estudiantes, familias/cuidadores e instituciones.
- Camila, 22: legalizaciones y fuente vigente.
- José, 34: formalización, secuencias, costos y tiempos.
- María, 48: gestiones familiares, simplicidad y conectividad intermitente.
- Early adopters: estudiantes UMSA, jóvenes profesionales y emprendedores de La Paz.

**Indicaciones visuales precisas**

Tres tarjetas de personas con inicial, edad, contexto, necesidad y valor. Debajo, matriz de segmentos y una banda de early adopters/canales.

**Instrucciones específicas para diagramas/tablas**

Fila superior: Camila, José y María. Fila inferior: cinco segmentos en dos filas. Resaltar instituciones en dorado para distinguir el componente B2B.

## 6. Modelo de negocio — Canvas

**Textos clave o viñetas**

- Freemium con información esencial gratuita y énfasis institucional.
- Premium: recordatorios avanzados, perfiles familiares, historial sincronizado y alertas.
- B2B: panel, estadísticas, configuración de flujos y white-label.
- No vender datos personales ni cobrar por información pública.

**Indicaciones visuales precisas**

Canvas completo de nueve bloques. Propuesta de valor en dorado; operación en azul; socios/segmentos en turquesa; relación/canales en violeta; costos en rojo e ingresos en verde.

**Instrucciones específicas para diagramas/tablas**

Distribución estándar: Socios | Actividades/Recursos | Propuesta | Relación/Canales | Segmentos. Abajo: Costos e Ingresos. Añadir una banda con el principio ético.

## 7. Tecnologías a utilizar

**Textos clave o viñetas**

- Flutter/Dart, Riverpod, Dio, Secure Storage y Drift.
- Go, Gin, REST/OpenAPI, pgx, sqlc y Goose en Render.
- PostgreSQL en Neon, full-text search y pgvector opcional.
- Gemini y Qwen mediante adaptadores, JSON Schema y validación.
- Docker, GitHub Actions, logs, métricas, límites y controles OWASP.

**Indicaciones visuales precisas**

Arquitectura en capas a la izquierda y pipeline RAG a la derecha; fondo azul marino.

**Instrucciones específicas para diagramas/tablas**

Capas: Clientes → API/Negocio → Datos → Operación. RAG: Consulta → Recuperación aprobada → Gemini/Qwen → Validación → Guía o abstención. Cerrar con controles de seguridad.

## 8. Prototipo o maquetado

**Textos clave o viñetas**

- Lenguaje simple, revelación progresiva, confianza y conectividad limitada.
- Pantallas: Inicio, Resultados, Ficha, Asistente, Checklist, Ubicación, Guardados, Reportar cambio y Perfil.
- La ficha muestra modalidad, fecha y fuente oficial.
- La prueba usa al menos diez participantes de tres perfiles y tres tareas.

**Indicaciones visuales precisas**

Tres maquetas de teléfono: Inicio, Ficha y Checklist. Panel lateral de principios UX y secuencia numerada inferior.

**Instrucciones específicas para diagramas/tablas**

En la ficha: modalidad, fecha, requisitos, costo/horario, ubicación y fuente. En el checklist: casillas y progreso. Flujo: Buscar → Verificar → Preparar → Derivar.

## 9. Análisis de mercado y competencia

**Textos clave o viñetas**

- Competidores/alternativas: gob.bo, Ciudadanía Digital/PTC, portales, buscadores/redes y tramitadores.
- Mi Trámite es complemento interoperable, no sustituto.
- Diferenciadores: personalización, trazabilidad, móvil, IA controlada y enfoque boliviano.
- Entrada: La Paz, 20–30 trámites, piloto universitario y demostración institucional.

**Indicaciones visuales precisas**

Tabla comparativa de seis filas con Mi Trámite resaltado; ruta de entrada al mercado debajo.

**Instrucciones específicas para diagramas/tablas**

Columnas: Alternativa | Fortaleza | Brecha/posición. Ruta: La Paz → 20–30 trámites → piloto universitario → demo institucional → expansión por paquetes.

## 10. Estudio económico

**Textos clave o viñetas**

- 97.350 Bs. de valor económico; 24.500 Bs. de desembolso.
- 6.000 Bs./mes de operación sostenible.
- Premium: 20 Bs./mes; licencia: 2.000 Bs./mes; white-label: 8.000–20.000 Bs.
- Equilibrio: 4 instituciones, 334 premium, o 2 instituciones + 134 premium.
- Escenario base: 3 instituciones y 250 premium hacia el último trimestre.

**Indicaciones visuales precisas**

Tres cifras principales; barras de inversión; barra apilada de operación; tarjetas de punto de equilibrio.

**Instrucciones específicas para diagramas/tablas**

Operación mensual: curación/soporte 2.500; IA 1.200; marketing 1.000; Render 400; contingencia 400; Neon 250; otros 250. Separar valor económico y desembolso.

## 11. Riesgos

**Textos clave o viñetas**

- Críticos: desactualización, alucinación, acceso no autorizado, confusión oficial, baja adopción y falta de actualización.
- Medios: caída de IA, costos variables y brecha digital.
- Alcance excesivo: probabilidad alta e impacto medio.
- Controles: versionado, fuentes, SLA, RAG, abstención, revisión, RBAC, límites, accesibilidad y aviso de independencia.

**Indicaciones visuales precisas**

Matriz 2×2 de probabilidad e impacto con diez riesgos; panel lateral con controles.

**Instrucciones específicas para diagramas/tablas**

Alta/media: alcance excesivo. Alta/alta: información desactualizada. Media/media: caída IA, costos y brecha digital. Media/alta: alucinación, acceso, confusión oficial, baja adopción y falta de actualización.

## 12. Cronograma

**Textos clave o viñetas**

- 20 semanas con sprints de dos semanas.
- Nueve fases desde descubrimiento hasta lanzamiento.
- Contenido y software avanzan en paralelo.
- Hitos: S3, S5, S11, S14, S17, S19 y S20.

**Indicaciones visuales precisas**

Diagrama de Gantt del documento casi a pantalla completa; banda inferior de hitos.

**Instrucciones específicas para diagramas/tablas**

Semanas 1–20 en columnas; fases en filas. Respetar: Descubrimiento S1–3; UX/UI S2–5; Catálogo S3–8; Backend S4–11; Flutter S5–13; IA/RAG S8–14; Integración/QA S12–17; Piloto S18–19; Lanzamiento S20.

## 13. Conclusiones y recomendaciones

**Textos clave o viñetas**

- Proyecto técnicamente alcanzable y socialmente relevante en aproximadamente cinco meses.
- Convierte información oficial en una guía personal, verificable y accionable.
- El principal activo es el catálogo confiable y actualizado.
- La empresa debe validar ingresos institucionales.
- Recomendaciones: limitar piloto, gobernar contenido, validar B2B y medir antes de migrar a Qwen local.

**Indicaciones visuales precisas**

Fondo azul marino; afirmación central; tres tarjetas de conclusiones; cuatro recomendaciones.

**Instrucciones específicas para diagramas/tablas**

Conclusiones: factibilidad técnica, activo de confianza y viabilidad híbrida. Recomendaciones numeradas abajo. Cierre: “menos incertidumbre + más preparación correcta + evidencia verificable”.
