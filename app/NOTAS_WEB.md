# Compatibilidad Web - Mi Trámite Bolivia

## ✅ Corrección Aplicada

La aplicación ahora funciona correctamente en web **sin anuncios**, ya que Google Mobile Ads no soporta la plataforma web.

## Cambios Implementados

### 1. Detección de Plataforma en `ad_config.dart`
```dart
import 'package:flutter/foundation.dart' show kIsWeb;

static bool get shouldShowAds {
  // No mostrar anuncios en web (no soportado)
  if (kIsWeb) {
    return false;
  }
  
  // Mostrar en móviles
  return true;
}
```

### 2. Inicialización Condicional en `main.dart`
```dart
void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Inicializar AdMob solo en plataformas móviles
  if (!kIsWeb) {
    await MobileAds.instance.initialize();
  }
  
  runApp(const MiTramiteApp());
}
```

## Comportamiento por Plataforma

| Plataforma | Logo Splash | Tema | Animaciones | Anuncios |
|------------|-------------|------|-------------|----------|
| **Android** | ✅ Sí | ✅ Sí | ✅ Sí | ✅ Sí |
| **iOS** | ✅ Sí | ✅ Sí | ✅ Sí | ✅ Sí |
| **Web** | ✅ Sí | ✅ Sí | ✅ Sí | ❌ No |

## Probar en Web

```bash
# Ejecutar en Chrome
flutter run -d chrome

# O en Edge
flutter run -d edge

# Construir para producción
flutter build web
```

## Alternativas de Monetización para Web

Si deseas monetizar la versión web, considera estas opciones:

### 1. Google AdSense (Recomendado para Web)
- Integración con HTML/JavaScript
- Usa paquete `dart:html` para inyectar scripts
- Ejemplo básico:
```dart
import 'dart:html' as html;

void initAdSense() {
  if (kIsWeb) {
    final script = html.ScriptElement()
      ..src = 'https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js'
      ..async = true;
    html.document.head?.append(script);
  }
}
```

### 2. Banner Placeholder (Actual)
La app muestra el diseño sin anuncios en web:
- Los espacios de anuncios simplemente no se renderizan
- `AdConfig.shouldShowAds` retorna `false` en web
- La experiencia de usuario es limpia y sin errores

### 3. Modelo Freemium
- Versión web gratuita sin anuncios
- Versión móvil con anuncios
- Suscripción premium para remover anuncios en móvil

## Ejecutar en Diferentes Plataformas

```bash
# Web (sin anuncios)
flutter run -d chrome

# Android (con anuncios de prueba)
flutter run -d android

# iOS (con anuncios de prueba)
flutter run -d ios

# Ver dispositivos disponibles
flutter devices
```

## Construcción para Producción

### Web
```bash
flutter build web --release
# Output en: build/web/
```

### Android
```bash
flutter build apk --release
# Output: build/app/outputs/flutter-apk/app-release.apk
```

### iOS
```bash
flutter build ios --release
# Luego usar Xcode para distribuir
```

## Nota sobre Rendimiento

La versión web tiene mejor rendimiento sin la inicialización de AdMob:
- ✅ Carga más rápida
- ✅ Menos consumo de memoria
- ✅ Sin dependencias nativas
- ✅ SEO-friendly (puede indexarse en Google)

## Recomendaciones

1. **Para demostración/testing**: Usa la versión web
2. **Para monetización con AdMob**: Usa versiones Android/iOS
3. **Para alcance máximo**: Publica en las 3 plataformas
   - Web: Sin anuncios, acceso universal
   - Android/iOS: Con anuncios, mayor monetización

---

**Fecha**: 27 de Julio, 2026
**Estado**: ✅ Funcional en todas las plataformas
