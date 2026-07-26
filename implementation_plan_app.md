# Mi Trámite Bolivia — App Ciudadana Flutter

Construir la aplicación ciudadana de **Mi Trámite Bolivia** como una aplicación móvil/cross-platform en **Flutter** con diseño moderno (Material 3), soporte offline y consumo de la API Go alineada con el documento de arquitectura del sistema (`mi-tramite-bolivia-sistema.md` v2.0) y requerimientos de app (`requerimientos_app.md`).

## Monetización y Suscripción Mensual (Demo)

Para la demostración del modelo de negocio de la app:
- **Pantalla dedicada `PremiumPlanScreen`**: Comparador de Plan Gratuito vs Plan Premium (Bs. 20/mes) con checkout simulado (pago QR / tarjeta demo).
- **Puntos de entrada de suscripción**:
  - Banner en `HomeScreen` ("Pasa a Premium para sincronizar tus checklists").
  - Botón de opción "Probar Premium" en `SavedScreen` y `ChecklistScreen`.
  - Opción destacada en `ProfileScreen` ("Mi Plan · Suscripción Activa / Cambiar Plan").
  - Badges explicativos de funciones premium (Alertas de cambios, Recordatorios avanzados, Perfiles familiares).

## Alignment with System Architecture (`mi-tramite-bolivia-sistema.md`)

- **Modelo Normalizado**: Soporte completo para instituciones, modalidades (presencial, en línea, mixta), requisitos condicionales (`aplica_si`), pasos numerados, costos desglosados, resultados y oficinas territoriales con horarios.
- **Trazabilidad y Verificación**: Despliegue visible de la fecha de última verificación (`verificado_en`), sello de verificación, aviso si la información requiere revisión y enlaces directos a las fuentes oficiales.
- **Privacidad y Uso Anónimo (APP-ONB-003, RF-CHAT-001)**: Sin obligatoriedad de registro ni solicitud de PII (carnet, correo, teléfono).
- **Asistente IA RAG con Citas (RF-CHAT-001 a 010)**: Consumo de la API de chat con renderizado de tarjetas de fuentes oficiales respaldadas y envío de feedback (👍/👎). Fallback estructurado si el LLM no tiene evidencia suficiente.
- **Reporte Ciudadano (RF-REP-001 a 005)**: Envío anónimo de reportes sobre trámites o datos desactualizados (`POST /api/v1/reportes`).

## Proposed Changes

Se inicializará el proyecto Flutter en `e:\UMSA\Emprendimiento\mi-tramite-bolivia\app`.

---

### Estructura del Proyecto Flutter

```
app/
├── pubspec.yaml                 ← Dependencias (http, shared_preferences, google_fonts, url_launcher, provider)
├── lib/
│   ├── main.dart                ← Entrypoint, MaterialApp, ThemeData y ruteo
│   ├── config/
│   │   ├── api_config.dart      ← Base URL (`http://localhost:8080` / `http://10.0.2.2:8080`), timeouts
│   │   └── theme.dart           ← Sistema de diseño Material 3 (Azul #2563EB, Turquesa #0F9F8F, Modo Oscuro)
│   ├── models/
│   │   ├── tramite.dart         ← Modelo completo de Trámite (modalidades, requisitos, pasos, costos, oficinas)
│   │   ├── categoria.dart       ← Modelo de Categoría
│   │   ├── institucion.dart     ← Modelo de Institución
│   │   ├── chat_message.dart    ← Modelo de Chat IA, Fuentes oficiales y Citas RAG
│   │   └── checklist_item.dart  ← Modelo de Requisito/Checklist local con notas
│   ├── services/
│   │   ├── api_service.dart     ← Wrapper HTTP para `/api/v1/*` (`tramites`, `categorias`, `instituciones`, `chat`, `reportes`)
│   │   ├── storage_service.dart ← Almacenamiento local (shared_preferences, favoritos, checklists, suscripción demo, caché offline)
│   │   └── connectivity_service.dart ← Detector de estado de red (online/offline banner)
│   ├── providers/
│   │   ├── app_provider.dart    ← Estado global (ciudad seleccionada, onboarding, estado offline, nivel de suscripción premium)
│   │   ├── favorites_provider.dart ← Gestión de trámites favoritos
│   │   ├── checklist_provider.dart ← Gestión de checklists interactivos y progreso
│   │   └── assistant_provider.dart ← Conversación RAG, estado del chat y feedback
│   ├── widgets/
│   │   ├── bottom_nav_bar.dart  ← Barra de navegación inferior (Inicio, Chat, Guardados, Perfil)
│   │   ├── tramite_card.dart    ← Card de trámite con badge de modalidad y sello de verificación
│   │   ├── category_chip.dart   ← Chip de categoría con icono
│   │   ├── verification_seal.dart ← Sello oficial con fecha `verificado_en`
│   │   ├── offline_banner.dart  ← Banner informativo de modo sin conexión
│   │   ├── disclaimer_banner.dart ← Aviso de independencia gubernamental
│   │   ├── requirement_tile.dart← Ítem de requisito marcable con checkbox y nota personal
│   │   ├── source_card.dart     ← Tarjeta de fuente oficial citada por IA
│   │   ├── cost_card.dart       ← Tarjeta con desglose de costo, concepto y modalidad
│   │   └── premium_cta_card.dart← Card/Banner promocional para suscripción mensual (Bs. 20/mes)
│   └── screens/
│       ├── onboarding_screen.dart ← Bienvenida (3 beneficios + aviso de independencia + ciudad)
│       ├── home_screen.dart       ← Inicio (Búsqueda prominente, categorías, frecuentes, checklists activos, CTA Premium)
│       ├── search_results_screen.dart ← Resultados de búsqueda con filtros (modalidad, institución, ciudad)
│       ├── tramite_detail_screen.dart ← Ficha completa del trámite (Requisitos, Pasos, Costos, Mapas, Botón Checklist)
│       ├── assistant_screen.dart     ← Chat RAG (burbujas, tarjetas de fuente, feedback 👍/👎, fallback)
│       ├── checklist_screen.dart     ← Gestor de checklists activos con barra de progreso animada
│       ├── saved_screen.dart         ← Trámites guardados, favoritos offline y sincronización premium
│       ├── profile_screen.dart       ← Configuración, ciudad, plan actual (Gratis/Premium), borrado de historial
│       ├── premium_plan_screen.dart  ← Comparador visual Plan Gratuito vs Premium (Bs. 20/mes) + Checkout demo (QR/Tarjeta)
│       └── report_form_screen.dart   ← Formulario de reporte ciudadano (datos obsoletos)
```

---

## Verification Plan

### Automated Verification
- `flutter pub get` para validar dependencias.
- `flutter analyze` para asegurar 0 errores de Dart/Flutter.

### Manual Verification
- Ejecutar la aplicación (`flutter run -d chrome` o emulador).
- Recorrido completo: Onboarding → Inicio → Búsqueda → Ficha de Trámite → Crear Checklist → Guardados → Asistente IA → Flujo de Suscripción Premium (Bs. 20/mes con QR/Simulación de Pago) → Reportar cambio.
- Simular desconexión a red para verificar operación en modo offline.
